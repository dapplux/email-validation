package emailvalidators

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	totalRequests = 50000
	parallelism   = 1000
	rateLimit     = 5000
)

type ZerobounceProvider struct {
	APIKey      string
	BaseURL     string
	Client      *http.Client
	RateLimiter *rate.Limiter
}

type ZerobounceValidationRequest struct {
	Email     string `json:"email"`
	APIKey    string `json:"api_key"`
	IPAddress string `json:"ip_address,omitempty"`
}

type ZerobounceValidationResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewZerobounceProvider(apiKey string) *ZerobounceProvider {
	return &ZerobounceProvider{
		APIKey:  apiKey,
		BaseURL: "https://api.zerobounce.net/v2/validate",
		Client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		RateLimiter: rate.NewLimiter(rate.Every(time.Second/rateLimit), rateLimit),
	}
}

func (zp *ZerobounceProvider) ValidateEmail(ctx context.Context, email string, ipAddress string) (*ZerobounceValidationResponse, error) {
	// Wait for the rate limiter
	if err := zp.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Build the query parameters
	params := url.Values{}
	params.Set("api_key", zp.APIKey)
	params.Set("email", email)
	if ipAddress != "" {
		params.Set("ip_address", ipAddress)
	}

	requestURL := fmt.Sprintf("%s?%s", zp.BaseURL, params.Encode())

	// Log the request URL
	//fmt.Printf("Request URL: %s\n", requestURL)

	// Make the GET request
	resp, err := zp.Client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var res ZerobounceValidationResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Log the response
	//fmt.Printf("Response dump: %+v\n", res)

	if res.Error != "" {
		return &res, fmt.Errorf("validation error: %s", res.Error)
	}

	return &res, nil
}

func (zp *ZerobounceProvider) BulkValidate(ctx context.Context, emails []string) {
	var wg sync.WaitGroup

	startTime := time.Now()
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := workerID; j < len(emails); j += parallelism {
				email := emails[j]
				response, err := zp.ValidateEmail(ctx, email, "")
				if err != nil {
					fmt.Printf("Error validating %s: %v\n", email, err)
				} else {
					fmt.Printf("Validated %s: %s\n", email, response.Status)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Completed %d validations in %v\n", len(emails), time.Since(startTime))
}
