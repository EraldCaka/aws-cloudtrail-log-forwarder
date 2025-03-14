package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookService struct {
	client *http.Client
}

func NewWebhookService() WebhookService {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return WebhookService{
		client: client,
	}
}

func (w *WebhookService) ForwardLog(ctx context.Context, callbackURL string, logData map[string]interface{}) error {
	logData["forwardedAt"] = time.Now().Format(time.RFC3339)

	jsonData, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("failed to marshal log data: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CloudTrailLogForwarder")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status code: %d", resp.StatusCode)
	}

	return nil
}

func (w *WebhookService) SendLog(callbackURL string, logData map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logData["forwardedAt"] = time.Now().Format(time.RFC3339)
	jsonData, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("failed to marshal log data: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CloudTrailLogForwarder")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status code: %d", resp.StatusCode)
	}

	return nil
}
