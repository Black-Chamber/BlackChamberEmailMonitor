package m365

import (
	"bcem/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// NewM365Instance initializes a new Microsoft 365 instance
func NewM365Instance(cfg config.Config) (*M365Instance, error) {

	return &M365Instance{
		TenantID:     cfg.Azure.TenantID,
		ClientID:     cfg.Azure.ClientID,
		ClientSecret: cfg.Azure.ClientSecret,
		APIMgmtURL:   "https://reports.office365.com/ecp/reportingwebservice/reporting.svc/MessageTrace",
	}, nil
}

// getOAuthToken obtains an OAuth2 token for the M365 instance
func (m *M365Instance) getOAuthToken() (*oauth2.Token, error) {
	config := clientcredentials.Config{
		ClientID:     m.ClientID,
		ClientSecret: m.ClientSecret,
		TokenURL: fmt.Sprintf(
			"https://login.microsoftonline.com/%s/oauth2/v2.0/token",
			m.TenantID,
		),
		Scopes: []string{"https://outlook.office365.com/.default"},
	}
	return config.Token(context.Background())
}

func (m *M365Instance) PerformLookup(startTime, endTime time.Time) (*MessageTraceResponse, error) {
	token, err := m.getOAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %v", err)
	}

	baseURL, err := url.Parse(m.APIMgmtURL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %v", err)
	}

	// Define query parameters
	filter := fmt.Sprintf("StartDate eq datetime'%s' and EndDate eq datetime'%s'",
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)
	queryParams := url.Values{}
	queryParams.Set("$filter", filter)

	var results []MessageTrace
	nextLink := baseURL.String() + "?" + queryParams.Encode()
	client := &http.Client{Timeout: 30 * time.Second}

	for nextLink != "" {
		req, err := http.NewRequest("GET", nextLink, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		req.Header.Set("Accept", "application/json")

		var resp *http.Response
		var retries int
		for retries < 5 {
			resp, err = client.Do(req)
			if err != nil {
				fmt.Printf("Backing off due to", err)
				retries++
				time.Sleep(time.Duration(retries) * time.Second) // Exponential backoff
				continue
			}
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				retries++
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			break
		}

		if err != nil || retries == 5 {
			return nil, fmt.Errorf("request failed after retries: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		// Parse response
		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var partialResp struct {
			Value    []MessageTrace `json:"value"`
			NextLink string         `json:"@odata.nextLink"`
		}
		if err := json.Unmarshal(responseBytes, &partialResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
		}

		results = append(results, partialResp.Value...)
		nextLink = partialResp.NextLink // Continue to the next page if available
	}

	return &MessageTraceResponse{Value: results}, nil
}
