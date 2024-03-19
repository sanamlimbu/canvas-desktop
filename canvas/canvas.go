package canvas

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ninja-software/terror/v2"
	"golang.org/x/time/rate"
)

type Controller struct {
	APIClient *APIClient
}

func NewController(client *APIClient) *Controller {
	return &Controller{
		APIClient: client,
	}
}

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
func (c *APIClient) do(req *http.Request) (*http.Response, error) {
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

func getNextURL(linkTxt string) string {
	url := ""
	if linkTxt != "" {
		links := strings.Split(linkTxt, ",")
		nextRegEx := regexp.MustCompile(`^<(.*)>; rel="next"$`)
		for i := 0; i < len(links); i++ {
			matches := nextRegEx.Match([]byte(links[i]))
			if matches {
				startIndex := strings.Index(links[i], "<")
				endIndex := strings.Index(links[i], ">")
				url = links[i][startIndex+1 : endIndex]
				break
			}
		}
	}

	return url
}

func (c *APIClient) GetAccessToken() string {
	accessToken := getenv("CANVAS_ACCESS_TOKEN", "")
	if accessToken != "" {
		return "error"
	}

	return accessToken
}

func getenv(key string, other string) string {
	value := os.Getenv(key)
	if value == "" {
		return other
	}

	return value
}

func (c *APIClient) LongRunningSleepFunc(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	for i := 0; i < 100; i++ {
		fmt.Println("index: ", i)
		time.Sleep(4 * time.Second)
	}
}
