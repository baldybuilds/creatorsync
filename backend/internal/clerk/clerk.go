package clerk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func Initialize() error {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return errors.New("CLERK_SECRET_KEY environment variable not set")
	}

	return nil
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		token := parts[1]

		// Log token info for debugging (only in development)
		if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "local" {
			fmt.Printf("Development mode: Token length: %d, JWT structure: %t\n",
				len(token),
				len(strings.Split(token, ".")) == 3)
		}

		if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "local" {
			// Development environment: Manual JWT parsing with relaxed validation
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token format - expected JWT with 3 parts",
				})
			}

			payload, err := decodeJWTSegment(parts[1])
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid token payload: %v", err),
				})
			}

			var claims map[string]interface{}
			if err := json.Unmarshal(payload, &claims); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid token claims: %v", err),
				})
			}

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or missing user ID in token",
				})
			}

			user := User{ID: userID}

			// Extract user info from claims
			if email, ok := claims["email"].(string); ok {
				user.Email = email
			}
			if firstName, ok := claims["first_name"].(string); ok {
				user.FirstName = firstName
			} else if firstName, ok := claims["firstName"].(string); ok {
				user.FirstName = firstName
			}
			if lastName, ok := claims["last_name"].(string); ok {
				user.LastName = lastName
			} else if lastName, ok := claims["lastName"].(string); ok {
				user.LastName = lastName
			}

			c.Locals("user", user)
			return c.Next()
		}

		// Production environment: Use Clerk SDK for strict validation
		sessionClaims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			// Provide more detailed error information
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Token verification failed: %v", err),
			})
		}

		userID := sessionClaims.Subject
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token claims",
			})
		}

		user := User{ID: userID}

		// Extract user info from verified claims
		var claimsMap map[string]interface{}
		if claimsBytes, err := json.Marshal(sessionClaims.Claims); err == nil {
			json.Unmarshal(claimsBytes, &claimsMap)

			if email, ok := claimsMap["email"].(string); ok {
				user.Email = email
			}
			if firstName, ok := claimsMap["first_name"].(string); ok {
				user.FirstName = firstName
			} else if firstName, ok := claimsMap["firstName"].(string); ok {
				user.FirstName = firstName
			}
			if lastName, ok := claimsMap["last_name"].(string); ok {
				user.LastName = lastName
			} else if lastName, ok := claimsMap["lastName"].(string); ok {
				user.LastName = lastName
			}
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func GetUserFromContext(c *fiber.Ctx) (*User, error) {
	user, ok := c.Locals("user").(User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, userID string) (*clerk.User, error) {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("CLERK_SECRET_KEY environment variable not set")
	}

	clerk.SetKey(secretKey)

	return user.Get(ctx, userID)
}

func decodeJWTSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	seg = strings.ReplaceAll(seg, "-", "+")
	seg = strings.ReplaceAll(seg, "_", "/")

	return base64.StdEncoding.DecodeString(seg)
}

func SyncUserData(ctx context.Context, userID string) error {
	// This is a placeholder for database sync logic
	// In a real implementation, you would:
	// 1. Get the user from Clerk
	// 2. Check if the user exists in your database
	// 3. Create or update the user in your database

	// For now, we'll just verify the user exists in Clerk
	_, err := GetUserByID(ctx, userID)
	return err
}
