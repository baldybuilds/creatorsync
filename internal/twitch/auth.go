package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ValidateToken(ctx context.Context, token string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute validation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var validationResp TokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
		return true, fmt.Errorf("failed to decode validation response: %w", err)
	}

	if validationResp.ClientID != "" {
		c.clientID = validationResp.ClientID
	}

	return true, nil
}
