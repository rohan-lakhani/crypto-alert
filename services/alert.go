package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"price-alert-system/models"

	"gopkg.in/mail.v2"
)

type AlertService struct {
	db               *sql.DB
	indicatorService *IndicatorService
	notificationChan chan models.Alert
	mutex            sync.Mutex
	smtpHost         string
	smtpPort         int
	smtpUsername     string
	smtpPassword     string
	fromEmail        string
}

func NewAlertService(db *sql.DB, indicatorService *IndicatorService, smtpHost string, smtpPort int, smtpUsername, smtpPassword, fromEmail string) *AlertService {
	return &AlertService{
		db:               db,
		indicatorService: indicatorService,
		notificationChan: make(chan models.Alert, 100),
		smtpHost:         smtpHost,
		smtpPort:         smtpPort,
		smtpUsername:     smtpUsername,
		smtpPassword:     smtpPassword,
		fromEmail:        fromEmail,
	}
}

func (s *AlertService) CreateAlert(alert *models.Alert) error {
	alert.Status = "pending"
	query := `INSERT INTO alerts (user_id, value, direction, indicator, status, email) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return s.db.QueryRow(query, alert.UserID, alert.Value, alert.Direction, alert.Indicator, alert.Status, alert.Email).Scan(&alert.ID)
}

func (s *AlertService) GetAlert(id int) (*models.Alert, error) {
	query := `SELECT id, user_id, value, direction, indicator, status, email FROM alerts WHERE id = $1`
	alert := &models.Alert{}
	err := s.db.QueryRow(query, id).Scan(&alert.ID, &alert.UserID, &alert.Value, &alert.Direction, &alert.Indicator, &alert.Status, &alert.Email)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *AlertService) GetPendingAlerts() ([]*models.Alert, error) {
	query := `SELECT id, user_id, value, direction, indicator, status, email FROM alerts WHERE status = 'pending' OR status = 'active'`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		alert := &models.Alert{}
		err := rows.Scan(&alert.ID, &alert.UserID, &alert.Value, &alert.Direction, &alert.Indicator, &alert.Status, &alert.Email)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func (s *AlertService) UpdateAlertStatus(id int, status string) error {
	query := `UPDATE alerts SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := s.db.Exec(query, status, id)
	return err
}

func (s *AlertService) CheckAlerts() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		alerts, err := s.GetPendingAlerts()
		if err != nil {
			log.Printf("Error fetching pending alerts: %v", err)
			continue
		}

		rsi, macd := s.indicatorService.GetIndicators()

		log.Printf("RSI: %f, MACD: %f", rsi, macd)
		log.Printf("Pending %d alerts", len(alerts))

		for _, alert := range alerts {
			go func(alert *models.Alert) {
				if alert.Status == "pending" {
					err := s.UpdateAlertStatus(alert.ID, "active")
					if err != nil {
						log.Printf("Error updating alert status to active: %v", err)
					}
					alert.Status = "active"
				}

				var currentValue float64
				switch strings.ToUpper(alert.Indicator) {
				case "RSI":
					currentValue = rsi
				case "MACD":
					currentValue = macd
				default:
					log.Printf("Unknown indicator: %s", alert.Indicator)
				}

				if (strings.ToUpper(alert.Direction) == "UP" && currentValue > alert.Value && currentValue > 0) ||
					(strings.ToUpper(alert.Direction) == "DOWN" && currentValue < alert.Value && currentValue > 0) {
					err := s.UpdateAlertStatus(alert.ID, "triggered")
					if err != nil {
						log.Printf("Error updating alert status to triggered: %v", err)
					}
					alert.Status = "triggered"
					s.notificationChan <- *alert

					// Send email notification
					if err := s.sendEmailNotification(alert, currentValue); err != nil {
						log.Printf("Error sending email notification: %v", err)
					}

					// After sending the notification, mark the alert as completed
					err = s.UpdateAlertStatus(alert.ID, "completed")
					log.Println("Alert completed:", alert.ID)
					if err != nil {
						log.Printf("Error updating alert status to completed: %v", err)
					}
				}
			}(alert)
		}
	}
}

func (s *AlertService) sendEmailNotification(alert *models.Alert, currentValue float64) error {
	// Fetch user email from the database based on alert.UserID
	var email string 
	err := s.db.QueryRow("SELECT email FROM alerts WHERE user_id = $1", alert.UserID).Scan(&email)
	if err != nil {
		return fmt.Errorf("error fetching user email: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Price Alert Triggered")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Price Alert Triggered</h1>
		<p>Your alert has been triggered:</p>
		<ul>
			<li>Indicator: %s</li>
			<li>Direction: %s</li>
			<li>Target Value: %.2f</li>
			<li>Current Value: %.2f</li>
		</ul>
	`, alert.Indicator, alert.Direction, alert.Value, currentValue))

	d := mail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(m)
}

func (s *AlertService) StartAlertChecker() {
	go s.CheckAlerts()
}

func (s *AlertService) GetNotificationChannel() <-chan models.Alert {
	return s.notificationChan
}