package httpx

import (
	"context"
	"io"
	"net/http"
	"time"
)

func DoRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	//Propagar headers (especialmente X-Request-ID)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 3 * time.Second, //Timeout general del cliente
	}

	return client.Do(req)
}
