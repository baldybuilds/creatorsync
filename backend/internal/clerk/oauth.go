package clerk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	clerk "github.com/clerk/clerk-sdk-go/v2"
)

func GetOAuthToken(ctx context.Context, userID, provider string) (string, error) {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("CLERK_SECRET_KEY environment variable not set")
	}

	clerk.SetKey(secretKey)

	user, err := GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	for _, account := range user.ExternalAccounts {
		if account.Provider == provider {
			url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s/oauth_access_tokens/%s", userID, provider)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return "", fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Add("Authorization", "Bearer "+secretKey)
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return "", fmt.Errorf("failed to make request: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read response body: %w", err)
			}

			var tokenResp struct {
				Data []struct {
					Token string `json:"token"`
				} `json:"data"`
			}

			if err := json.Unmarshal(bodyBytes, &tokenResp); err == nil {
				if len(tokenResp.Data) > 0 {
					return tokenResp.Data[0].Token, nil
				}
			}

			var tokensArray []struct {
				Token string `json:"token"`
			}

			if err := json.Unmarshal(bodyBytes, &tokensArray); err == nil {
				if len(tokensArray) > 0 {
					return tokensArray[0].Token, nil
				}
			}

			var tokenObj struct {
				Token string `json:"token"`
			}

			if err := json.Unmarshal(bodyBytes, &tokenObj); err == nil {
				if tokenObj.Token != "" {
					return tokenObj.Token, nil
				}
			}
			return "", fmt.Errorf("could not parse token from response: %s", string(bodyBytes))
		}
	}

	return "", fmt.Errorf("user does not have a connected %s account", provider)
}
