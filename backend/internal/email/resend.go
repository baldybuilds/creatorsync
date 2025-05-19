package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ResendClient struct {
	apiKey     string
	apiBaseURL string
}
type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}
type WaitlistRequest struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}
type EmailResponse struct {
	ID    string `json:"id"`
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

func NewResendClient() (*ResendClient, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, errors.New("RESEND_API_KEY environment variable is not set")
	}

	return &ResendClient{
		apiKey:     apiKey,
		apiBaseURL: "https://api.resend.com",
	}, nil
}

func (c *ResendClient) AddToWaitlist(req WaitlistRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailReq := EmailRequest{
		From:    "waitlist@creatorsync.app",
		To:      []string{req.Email},
		Subject: "Welcome to CreatorSync Waitlist!",
		HTML: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
				<h1 style="color: #6366f1;">Welcome to CreatorSync!</h1>
				<p>Hi %s,</p>
				<p>Thank you for joining our waitlist! We're excited to have you on board.</p>
				<p>We're working hard to build the best platform for creators to streamline their content workflow.</p>
				<p>We'll notify you as soon as we're ready to welcome you to our beta program.</p>
				<p>Best regards,<br>The CreatorSync Team</p>
			</div>
		`, req.Name),
	}

	adminEmailReq := EmailRequest{
		From:    "waitlist@creatorsync.app",
		To:      []string{"info@creatorsync.app"},
		Subject: "New Waitlist Signup",
		HTML: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
				<h1 style="color: #6366f1;">New Waitlist Signup</h1>
				<p>A new user has joined the waitlist:</p>
				<p><strong>Email:</strong> %s</p>
				<p><strong>Name:</strong> %s</p>
			</div>
		`, req.Email, req.Name),
	}

	if err := c.sendEmail(emailReq); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	if err := c.sendEmail(adminEmailReq); err != nil {
		fmt.Printf("Failed to send admin notification: %v\n", err)
	}

	return nil
}

func (c *ResendClient) sendEmail(req EmailRequest) error {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.apiBaseURL+"/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if len(respBody) > 0 {
		var emailResp EmailResponse
		if err := json.Unmarshal(respBody, &emailResp); err != nil {
			return fmt.Errorf("failed to decode response: %w, body: %s", err, string(respBody))
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("resend API error: %s (code: %s)", emailResp.Error.Message, emailResp.Error.Code)
		}
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend API error: status code %d with empty response", resp.StatusCode)
	}

	return nil
}
