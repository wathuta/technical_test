package model

type TransactionType string

const (
	CustomerPayBillOnline TransactionType = "CustomerPayBillOnline"
)

type Currency string

const (
	KES Currency = "KES"
)
const BaseURL = "https://sandbox.safaricom.co.ke"
const BusinessSortCode = "174379"

type STKPushRequestBody struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// STKPushRequestResponse is the response sent back after initiating an STK push request.
type STKPushRequestResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
	RequestID           string `json:"requestId"`
	ErrorCode           string `json:"errorCode"`
	ErrorMessage        string `json:"errorMessage"`
}

// MpesaAccessTokenResponse is the response sent back by Safaricom when we make a request to generate a token
type MpesaAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RequestID    string `json:"requestId"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type Item struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}

type CallbackMetadata struct {
	Item []Item `json:"Item"`
}

type StkCallback struct {
	MerchantRequestID string           `json:"MerchantRequestID"`
	CheckoutRequestID string           `json:"CheckoutRequestID"`
	ResultCode        int              `json:"ResultCode"`
	ResultDesc        string           `json:"ResultDesc"`
	CallbackMetadata  CallbackMetadata `json:"CallbackMetadata"`
}

type Body struct {
	StkCallback StkCallback `json:"stkCallback"`
}

type CallbackResponse struct {
	Body Body `json:"Body"`
}
