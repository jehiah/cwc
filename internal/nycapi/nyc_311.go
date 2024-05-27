// https://api-portal.nyc.gov/profile
// for NYC 311 Public Developers
package nycapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ServiceRequest struct {
	SRNumber                    string `json:"SRNumber"`
	Agency                      string `json:"Agency"`
	Problem                     string `json:"Problem"`
	ProblemDetails              string `json:"ProblemDetails"`
	ResolutionActionUpdatedDate string `json:"ResolutionActionUpdatedDate"`
	Status                      string `json:"Status"`
	DateTimeSubmitted           string `json:"DateTimeSubmitted"`
	ResolutionAction            string `json:"ResolutionAction"`
	Address                     struct {
		Borough     string `json:"Borough"`
		FullAddress string `json:"FullAddress"`
	} `json:"Address"`
}

// Client is a client for the NYC Public API
type Client struct {
	SubscriptionKey string
}

func (c *Client) GetServiceRequest(ctx context.Context, srNumber string) (*ServiceRequest, error) {
	// GET https://api.nyc.gov/public/api/GetServiceRequest?srnumber=311-18246322 HTTP/1.1
	// Ocp-Apim-Subscription-Key:

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.nyc.gov/public/api/GetServiceRequest?srnumber="+url.QueryEscape(srNumber), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", c.SubscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var data ServiceRequest
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
