
# AWS CloudTrail Log Forwarder

## Description

The **AWS CloudTrail Log Forwarder** is a service designed to fetch and forward AWS CloudTrail logs to a specified webhook URL.
It integrates with AWS services, MongoDB, Redis, and provides real-time log forwarding capabilities.
The service also allows for managing data sources and is built with the Fiber framework for fast and scalable web service handling.
(Also has mock data if you want to test the application without AWS CloudTrail real logs.)

- Will add monitoring and logging features with Prometheus and Grafana.
- Will add recovery and retry mechanisms for failed log forwarding attempts.
- Will implement rate limiting to prevent abuse and ensure fair usage.
- Will implement data encryption and secure storage for sensitive information.

### Key Features:
- Fetches AWS CloudTrail logs for specific time intervals.
- Forwards logs to a configured webhook endpoint.
- Integrates with MongoDB for storing logs and metadata.
- Uses Redis for task queuing and caching.
- Provides an HTTP API for managing data sources (add, remove, list sources).
- Real-time log fetching with configurable intervals.

## Prerequisites

Before running the application, ensure you have the following:

- **Docker** (for containerized services like MongoDB and Redis)
- **Go** (for building and running the application)
- **AWS** account with CloudTrail enabled and logs configured.
- **MongoDB** (for storing log data and sources)
- **Redis** (for caching and task queuing)

## Getting Started

Follow these steps to get the project up and running locally:

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/aws-cloudtrail-log-forwarder.git
cd aws-cloudtrail-log-forwarder
```

### 2. Set up MongoDB and Redis (Using Docker Compose)

If you don't have MongoDB and Redis running, you can use the provided `docker-compose.yml` to set them up:

```bash
docker-compose up -d
```

This will start MongoDB and Redis services in detached mode.

### 3. Configure MongoDB and Redis

Ensure MongoDB is running on `localhost:27017` and Redis is running on `localhost:6379`. You may need to adjust connection strings in the code if your environment differs.

### 4. Set Up AWS Credentials

Make sure you have your AWS and MongoDB credentials configured on `.env` file.

```bash
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
MONGODBCONNECTION=mongodb://root:rootpassword@localhost:27017
```

### 5. Accessing the API

You can interact with the API using HTTP requests. The following routes are available:

- **GET `/sources`**: Fetch all configured sources.
- **POST `/sources`**: Add a new source.
- **DELETE `/sources/:id`**: Delete a source by its ID.

### 6. Real-Time Log Fetching and Forwarding

The application fetches CloudTrail logs in real-time at a configurable interval every minute. These logs are forwarded to the configured webhook endpoint for further processing.

### 7. Example Source JSON

When adding a source using the `POST /sources` endpoint, the following JSON format is used:

```json
{
  "sourceType": "cloudtrail",
  "accessKeyId": "test",
  "secretAccessKey": "test",
  "region": "us-west-2",
  "logFetchInterval": 15,
  "callbackUrl": "http://localhost:8080/webhook",
  "s3Bucket": "logsBuckTest",
  "s3Prefix": "AWSLogs/"
}
```

### 8. Monitoring Logs

You can monitor the logs in the terminal or configure them to be stored in MongoDB for later retrieval. Logs are stored with metadata like `eventId`, `eventName`, and `eventTime` for later analysis.



