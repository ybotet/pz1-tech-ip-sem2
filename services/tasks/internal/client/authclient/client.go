package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ybotet/pz1-tech-ip-sem2/shared/httpx"
)

type AuthClient struct {
	baseURL string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{baseURL: baseURL}
}

// VerifyToken llama al endpoint /v1/auth/verify del servicio Auth
func (c *AuthClient) VerifyToken(ctx context.Context, token string, requestID string) (bool, string, error) {
	url := fmt.Sprintf("%s/v1/auth/verify", c.baseURL)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"X-Request-ID":  requestID,
	}

	// Usar el cliente compartido con timeout
	resp, err := httpx.DoRequest(ctx, http.MethodGet, url, nil, headers)
	if err != nil {
		// Timeout o error de conexión
		return false, "", fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Valid   bool   `json:"valid"`
		Subject string `json:"subject"`
		Error   string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return result.Valid, result.Subject, nil
	case http.StatusUnauthorized:
		return false, "", nil // Token inválido, no es un error del sistema
	default:
		return false, "", fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}
}
