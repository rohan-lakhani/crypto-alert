package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"price-alert-system/models"
	"price-alert-system/services"
)

func RegisterRoutes(e *echo.Echo, alertService *services.AlertService) {
	e.POST("/alerts", createAlert(alertService))
	e.GET("/alerts/:id", getAlert(alertService))
}

func createAlert(alertService *services.AlertService) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqAlert := new(models.AlertRequest)
		if err := c.Bind(reqAlert); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid alert data"})
		}

		log.Println(reqAlert)

		alert := &models.Alert{
			UserID:    reqAlert.UserID,
			Value:     reqAlert.Value,
			Direction: reqAlert.Direction,
			Indicator: reqAlert.Indicator,
			Email:     reqAlert.Email,
			Status:    "pending",
		}
		err := alertService.CreateAlert(alert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create alert"})
		}

		return c.JSON(http.StatusCreated, alert)
	}
}

func getAlert(alertService *services.AlertService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid alert ID"})
		}

		alert, err := alertService.GetAlert(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch alert"})
		}

		return c.JSON(http.StatusOK, alert)
	}
}
