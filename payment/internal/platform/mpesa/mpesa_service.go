package mpesa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wathuta/technical_test/payment/internal/model"
)

type MpesaService interface {
	InitiateSTKPushRequest(body *model.STKPushRequestBody) (*model.STKPushRequestResponse, error)
}

// Mpesa is an application that will be making a transaction
type Mpesa struct {
	consumerKey    string
	consumerSecret string
	baseURL        string
	client         *http.Client
}

// MpesaOpts stores all the configuration keys we need to set up a Mpesa app,
type MpesaOpts struct {
	ConsumerKey    string
	ConsumerSecret string
}

// NewMpesa sets up and returns an instance of Mpesa
func NewMpesa(m *MpesaOpts) MpesaService {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Mpesa{
		consumerKey:    m.ConsumerKey,
		consumerSecret: m.ConsumerSecret,
		baseURL:        model.BaseURL,
		client:         client,
	}
}

// makeRequest performs all the http requests for the specific app
func (m *Mpesa) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// initiateSTKPushRequest makes a http request performing an STK push request
func (m *Mpesa) InitiateSTKPushRequest(body *model.STKPushRequestBody) (*model.STKPushRequestResponse, error) {
	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", m.baseURL)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stk push request json with error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create stk push request with error: %w", err)
	}

	accessTokenResponse, err := m.generateAccessToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessTokenResponse.AccessToken))

	resp, err := m.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("stk push request failed with error: %w", err)
	}

	stkPushResponse := new(model.STKPushRequestResponse)
	if err := json.Unmarshal(resp, &stkPushResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stkpush response error: %w", err)
	}

	if stkPushResponse.ErrorCode != "0" && stkPushResponse.ErrorCode != "" {
		return nil, fmt.Errorf("initiate stkpush request failed with error: %s", stkPushResponse.ErrorMessage)
	}

	return stkPushResponse, nil
}

// generateAccessToken sends a http request to generate new access token
func (m *Mpesa) generateAccessToken() (*model.MpesaAccessTokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", m.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ucreate generate access token request with error: %w", err)
	}

	req.SetBasicAuth(m.consumerKey, m.consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("generate access token request failed with error: %w", err)
	}

	accessTokenResponse := new(model.MpesaAccessTokenResponse)
	if err := json.Unmarshal(resp, &accessTokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal generate access token response error: %w", err)
	}

	return accessTokenResponse, nil
}
