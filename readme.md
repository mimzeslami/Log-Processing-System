# Log Processing System

This project implements a log processing system using RabbitMQ and PostgreSQL. It consists of two main components:
1. **Log Generator**: Simulates structured log events and pushes them to RabbitMQ.
2. **Log Processor**: Consumes logs from RabbitMQ, processes them, and stores the structured logs in PostgreSQL.

---

## Features
- Simulates structured logs in the format:
  ```
  status_code,api_endpoint,message,timestamp,ip_address
  ```
- Processes logs with validation and parsing.
- Stores structured logs in a PostgreSQL database for analysis.
- Uses RabbitMQ for reliable message queuing.

---

## Project Structure

```
log-processing-system/
├── cmd/
│   ├── log-generator/        # Log generator service
│   │   └── main.go
│   ├── log-processor/        # Log processor service
│   │   └── main.go
├── internal/
│   ├── db/                   # Database connection and operations
│   │   └── db.go
│   ├── mq/                   # RabbitMQ integration
│   │   └── rabbitmq.go
│   ├── logprocessor/         # Log processing logic
│   │   └── processor.go
├── config/
│   └── config.go             # Configuration management
├── docker-compose.yml        # Docker setup for RabbitMQ and PostgreSQL
├── .env                      # Environment variables (not included in Git)
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency lock file
└── README.md                 # Project documentation
```

---

## Prerequisites
- [Go](https://go.dev/) (>= 1.17)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/log-processing-system.git
cd log-processing-system
```

### 2. Configure Environment Variables
Create a `.env` file in the project root with the following content:
```env
DB_CONN_STRING=postgres://loguser:logpassword@localhost:5432/logprocessor?sslmode=disable
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=logs_queue
```

### 3. Start Services with Docker Compose
```bash
docker-compose up -d
```

This will start:
- **PostgreSQL** on port `5432`
- **RabbitMQ** on ports `5672` (AMQP) and `15672` (Management UI)

### 4. Run the Log Generator
The log generator simulates structured logs and sends them to RabbitMQ.
```bash
go run cmd/log-generator/main.go
```

### 5. Run the Log Processor
The log processor consumes logs from RabbitMQ, processes them, and stores them in the PostgreSQL database.
```bash
go run cmd/log-processor/main.go
```

---

## Log Template

Logs are structured as:
```
status_code,api_endpoint,message,timestamp,ip_address
```

**Example Logs**:
```
200,/api/login,User logged in,2024-12-21T12:00:00Z,192.168.1.10
500,/api/upload,Server error occurred,2024-12-21T12:01:00Z,192.168.1.11
403,/api/profile,Access denied,2024-12-21T12:02:00Z,192.168.1.12
```

---

## Database Schema

Create the following table in PostgreSQL to store structured logs:

```sql
CREATE TABLE structured_logs (
    id SERIAL PRIMARY KEY,
    status_code INT NOT NULL,
    api_endpoint TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    ip_address TEXT NOT NULL
);
```

---

## Testing the System

1. Start the RabbitMQ Management UI:
   - URL: [http://localhost:15672](http://localhost:15672)
   - Username: `guest`
   - Password: `guest`

2. Monitor the `logs_queue` in the RabbitMQ Management UI for incoming messages.

3. Query the database to view processed logs:
```sql
SELECT * FROM structured_logs ORDER BY timestamp DESC;
```

---

## Example Workflow

1. **Log Generator Output**:
```
Published log: 200,/api/login,User logged in,2024-12-21T12:00:00Z,192.168.1.10
Published log: 500,/api/upload,Server error occurred,2024-12-21T12:01:00Z,192.168.1.11
```

2. **Log Processor Output**:
```
Processing log: 200,/api/login,User logged in,2024-12-21T12:00:00Z,192.168.1.10
Successfully processed log: {200 /api/login User logged in 2024-12-21 12:00:00 +0000 UTC 192.168.1.10}
Processing log: 500,/api/upload,Server error occurred,2024-12-21T12:01:00Z,192.168.1.11
Successfully processed log: {500 /api/upload Server error occurred 2024-12-21 12:01:00 +0000 UTC 192.168.1.11}
```

3. **Database Query Output**:
```text
 id | status_code | api_endpoint |         message         |       timestamp       |  ip_address  
----+-------------+--------------+-------------------------+-----------------------+--------------
  1 |         200 | /api/login   | User logged in          | 2024-12-21 12:00:00  | 192.168.1.10
  2 |         500 | /api/upload  | Server error occurred   | 2024-12-21 12:01:00  | 192.168.1.11
```

---

## To-Do and Enhancements

- [ ] Add retries for failed message processing.
- [ ] Configure a dead-letter queue for unprocessable messages.
- [ ] Implement a monitoring tool to track message processing status.
- [ ] Deploy using Docker for the entire stack (including the application).

---

## License
This project is licensed under the MIT License. Feel free to use, modify, and distribute it. 😊
