# **AWS CloudTrail Log Forwarder – Backend Challenge**

## **Introduction**
Welcome to the **Cybee.ai Backend Challenge**! 🎯

Your task is to build a **high-performance system** that integrates with **AWS CloudTrail**, fetches event logs, and forwards them to a specified webhook. This challenge will assess your ability to:
✅ **Build a secure REST API** (Fast API with Golang)
✅ **Fetch logs from AWS CloudTrail** (via AWS SDK)
✅ **Schedule background jobs** (using Redis-based queues)
✅ **Handle failures, rate limits, and retries**
✅ **Store credentials securely**

Successful completion of this challenge will earn you a **technical interview with our team!** 🚀

---

## **Tech Stack Requirements**

Your solution **must** use:
- **Golang** (for API development)
- **Fiber** or **Echo** (for REST API)
- **MongoDB** (for storing sources and logs)
- **Redis + Asynq** (for job scheduling and retries)
- **AWS SDK (Go v2)** (for fetching CloudTrail logs)

**Bonus (optional, but a plus! 💡)**
- **Elasticsearch** (for log indexing and querying)
- **Docker & Deployment** (containerized setup)

---

## **Project Requirements**

### **1. REST API (Fast, Secure, and Scalable 🚀)**

Implement the following endpoints using **Fiber** or **Echo**:

#### **🔹 POST `/add-source`**
Registers a new AWS CloudTrail source and stores credentials securely.

#### **🔹 DELETE `/remove-source/:id`**
Removes an existing event source.

#### **🔹 GET `/sources`**
Returns a list of active sources.

**Data Model (MongoDB)**
```json
{
  "id": "uuid",
  "sourceType": "aws_cloudtrail",
  "credentials": {
    "accessKeyId": "string",
    "secretAccessKey": "string",
    "region": "us-east-1"
  },
  "logFetchInterval": 300,
  "callbackUrl": "https://example.com/webhook",
  "s3Bucket": "cloudtrail-logs",
  "s3Prefix": "AWSLogs/123456789012/CloudTrail/"
}
```
2. Fetch & Forward Logs Automatically

Once a source is added:

✅ Schedule a background job to fetch logs periodically (Redis + Asynq)

✅ Use AWS SDK to fetch logs from CloudTrail

✅ Forward logs to the callbackUrl

✅ Handle retries & failures (e.g., webhook down, AWS rate limits)


Example CloudTrail Log

```json
{
  "eventId": "123456",
  "eventTime": "2024-03-12T14:00:00Z",
  "eventType": "AWSConsoleSignIn",
  "sourceIPAddress": "192.168.1.1",
  "userIdentity": {
    "arn": "arn:aws:iam::123456789012:user/admin"
  }
}
```
3. Handling Edge Cases & Failure Scenarios

Your system must handle:

🚦 API Rate Limits → Use exponential backoff + retries

🔑 Credential Expiry → Detect and notify the user

🔁 Duplicate Logs → Ensure logs aren’t sent twice

💥 Webhook Failures → Retry failed webhook deliveries
```bash
go get github.com/gofiber/fiber/v2
go get github.com/aws/aws-sdk-go-v2
go get go.mongodb.org/mongo-driver/mongo
go get github.com/hibiken/asynq
go get github.com/go-resty/resty/v2
```

Evaluation Criteria

✅ Code Quality – Clean, modular, well-documented

✅ Performance – Efficient API and log processing

✅ Error Handling – Robust handling of edge cases

✅ Security – Secure storage of credentials

✅ Scalability – Handles high log volume efficiently

🚀 Bonus Points (Not Required, But Impressive! 🌟)


✅ Dockerized Deployment

✅ Logging & Monitoring (e.g., Prometheus, ELK Stack)

✅ Unit Tests (for critical components)

✅ CI/CD Integration
