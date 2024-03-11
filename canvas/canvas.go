package canvas

import (
	"context"
	"net/http"

	"github.com/ninja-software/terror/v2"
	"golang.org/x/time/rate"
)

type APIClient struct {
	BaseURL      string
	AccessToken  string
	PageSize     int
	Client       *http.Client
	RateLimitter *rate.Limiter
}

func NewAPIClient(baseURL string, accessToken string, pageSize int, client *http.Client, rateLimitter *rate.Limiter) *APIClient {
	return &APIClient{
		BaseURL:      baseURL,
		AccessToken:  accessToken,
		PageSize:     pageSize,
		Client:       client,
		RateLimitter: rateLimitter,
	}
}

// https://medium.com/mflow/rate-limiting-in-golang-http-client-a22fba15861a
func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.RateLimitter.Wait(ctx)
	if err != nil {
		return nil, terror.Error(err, "error rate limmiter wait")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, terror.Error(err, "error sending HTTP request")
	}

	return resp, nil
}
