package facebook

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
)

const (
	devAccessToken 		string = "1054196554605437|4Lm18mcDz88yNZzOuznjAcQo4bI"
	stagingAccessToken 	string = "1051736634851429|hfSuhz9DrbjcIKF2vOZLV_ooNrk"
	facebookUrl 		string = "https://graph.facebook.com"
)

//Error basic struct
type Error struct {
	error
}

// Facebook receipt
type Receipt struct {
	Id    string    `json:"id"`
	User  User	`json:"user"`
	Application Application `json:"application"`
	Actions []Action `json:"actions"`
	RefundableAmount RefundableAmount `json:"refundable_amount"`
	Items []Item `json:"items"`
	Country string `json:"country"`
	CreatedTime string `json:"created_time"`
	PayoutForeignExchangeRate float32 `json:"payout_foreign_exchange_rate"`
}

// User struct
type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// Application structure
type Application struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}

// Action structure
type Action struct {
	Type string `json:"type"`
	Status string `json:"status"`
	Currency string `json:"currency"`
	Amount string `json:"amount"`
	TimeCreated string `json:"time_created"`
	TimeUpdated string `json:"time_updated"`
}

// Refundable amount structure
type RefundableAmount struct {
	Currency string `json:"currency"`
	Amount string `json:"amount"`
}

// Item structure
type Item struct {
	Type string `json:"type"`
	Product string `json:"product"`
	Quantity int `json:"quantity"`
}

// Facebook error structure
type FacebookError struct {
	Error ErrorMessage `json:"error"`
}

// Facebook error message
type ErrorMessage struct {
	Message string `json:"message"`
	Type string `json:"type"`
	Code int `json:"code"`
	FBTraceId string `json:"fbtrace_id"`
}

//VerifyReceipt will check for the given facebook payment id
func VerifyReceiptByPaymentId(paymentID string) (*Receipt, error) {
	receipt, err := sendReceiptToFacebook(paymentVerificationURL(devAccessToken, paymentID))
	return receipt, err
}

// Build payment verification url
func paymentVerificationURL(accessToken string, paymentID string) string {
	return facebookUrl + "/" + paymentID + "?access_token=" + accessToken
}

// Sends the receipt to facebook, returns the receipt or an error upon completion
func sendReceiptToFacebook(url string) (*Receipt, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		facebookError := new(FacebookError)
		_ = json.Unmarshal(body, &facebookError)

		return nil, verificationError(resp.StatusCode, facebookError)
	}

	receipt := new(Receipt)
	err = json.Unmarshal(body, &receipt)

	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// Generates the correct error based on a status error code
func verificationError(errCode int, facebookError *FacebookError) error {
	var errorMessage string

	switch errCode {
	case 400:
		errorMessage = facebookError.Error.Message
		break
	case 404:
		errorMessage = facebookError.Error.Message
		break
	case 500:
		errorMessage = "There was an Internal Server Error."
		break
	default:
		errorMessage = "An unknown error ocurred."
		break
	}

	return &Error{errors.New(errorMessage)}
}