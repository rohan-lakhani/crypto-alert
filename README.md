# Alert System

This project is an alert system that allows users to create and manage alerts based on specific indicators and values.

## Getting Started

Follow these steps to set up and run the project:

1. Clone the repository:
   ```
   git clone https://github.com/rohan-lakhani/crypto-alert.git
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run main.go
   ```

   After running this command, you will see all pending alerts in the log.

## Usage

### Creating an Alert

To create a new alert, send a POST request to `http://localhost:3030/alerts` using Postman or any other API client.

Request body example:
```json
{
    "user_id": 1,
    "email": "rohanlakhani2003@gmail.com",
    "value": 2,
    "direction": "UP",
    "indicator": "RSI"
}
```

Response example:
```json
{
    "id": 71,
    "user_id": 1,
    "email": "rohanlakhani2003@gmail.com",
    "value": 2,
    "direction": "UP",
    "indicator": "RSI",
    "status": "pending"
}
```

### Retrieving a Specific Alert

To get information about a specific alert, send a GET request to `http://localhost:3030/alerts/:alertId`, where `:alertId` is the ID of the alert you want to retrieve.

Response example:
```json
{
    "id": 71,
    "user_id": 1,
    "email": "rohanlakhani2003@gmail.com",
    "value": 2,
    "direction": "UP",
    "indicator": "RSI",
    "status": "active"
}
```

## Alert Notifications

When an alert is triggered, you will receive an email at the address you provided when creating the alert.


## Contact

if you have any questions reach me at rohanlakhani2003@gmail.com