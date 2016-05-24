package amazon

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"errors"
)

const (
	developerSecret     string = "2:smXBjZkWCxDMSBvQ8HBGsUS1PK3jvVc8tuTjLNfPHfYAga6WaDzXJPoWpfemXaHg:iEzHzPjJ-XwRdZ4b4e7Hxw=="
	amazonSandboxURL    string = "http://localhost:8080/RVSSandbox/"
	amazonProductionURL string = "https://appstore-sdk.amazon.com/version/1.0/verifyReceiptId/"
)

type Error struct {
	error
}

// Receipt is information returned by Amazon
// Documentation: https://developer.amazon.com/public/apis/earn/in-app-purchasing/docs-v2/verifying-receipts-in-iap-2.0
type Receipt struct {
	PurchaseDate        int            `json:"purchaseDate"`
	RenewalDate         string            `json:"renewalDate"`
	ReceiptID           string            `json:"receiptID"`
	ProductID           string            `json:"productID"`
	ParentProductID     string            `json:"parentProductID"`
	ProductType         string            `json:"productType"`
	CancelDate          string            `json:"cancelDate"`
	Term                string            `json:"term"`
	TermSku             string            `json:"termSku"`
	Quantity            int               `json:"quantity"`
	BetaProduct         bool              `json:"betaProduct"`
	TestTransaction     bool              `json:"testTransaction"`
}

// For the given amazon userId and receiptId verify
func VerifyReceipt(userId string, receiptId string,  useSandbox bool) (*Receipt, error) {
	receipt, err := sendReceiptToAmazon(userId, receiptId, verificationURL(useSandbox))
	return receipt, err
}

// Selects the proper url to use when talking to apple based on if we should use the sandbox environment or not
func verificationURL(useSandbox bool) string {

	if useSandbox {
		return amazonSandboxURL
	}
	return amazonProductionURL
}

// Build final url that we will call
func buildFinalUrl(url, userId, receiptId string) string {
	return url + "developer/" + developerSecret + "/user/" + userId + "/receiptId/" + receiptId
}

// Sends the receipt to apple, returns the receipt or an error upon completion
func sendReceiptToAmazon(userId, receiptId, url string) (*Receipt, error) {

	resp, err := http.Get(buildFinalUrl(url, userId, receiptId))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	receipt := new(Receipt)
	err = json.Unmarshal(body, &receipt)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, verificationError(resp.StatusCode)
	}

	return receipt, nil
}

// Generates the correct error based on a status error code
func verificationError(errCode int) error {
	var errorMessage string

	switch errCode {
		case 400:
			errorMessage = "The transaction represented by this receiptId is invalid, or no transaction was found for this receiptId."
			break
		case 496:
			errorMessage = "Invalid sharedSecret"
			break

		case 497:
			errorMessage = "Invalid User ID"
			break

		case 500:
			errorMessage = "There was an Internal Server Error"
			break
		default:
			errorMessage = "An unknown error ocurred"
			break
	}

	return &Error{errors.New(errorMessage)}
}
