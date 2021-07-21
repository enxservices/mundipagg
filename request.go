package mundipagg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BASEURL                             string = "https://api.mundipagg.com/core/v1/"
	CUSTOMERURL                         string = BASEURL + "customers"
	SUBSCRIPTIONURL                     string = BASEURL + "subscriptions"
	SUBSCRIPTIONUPDATENEXTBILLINGDAYURL string = "/billing-date"
	SUBSCRIPTIONUPDATESTARTATURL        string = "/start-at"
	CARDENDPOINT                        string = "cards"
	DISCOUNTENDPOINT                    string = "/discounts"
)

type Response struct {
	ID         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Email      string     `json:"email,omitempty"`
	Delinquent bool       `json:"delinquent,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`

	MundipaggJSONAnswer string
}

func Do(method string, data interface{}, secretKey string, indepotencyKey string, url string) (*Response, error) {
	postData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString([]byte(secretKey + ":"))
	req.Header.Set("Authorization", "Basic "+token)
	req.Header.Set("Content-Type", "application/json")
	if indepotencyKey != "" {
		req.Header.Set("Idempotency-Key", indepotencyKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("Invalid Request:\nSent:\n" + string(postData) + "Received:\n" + string(body))
	}

	response := &Response{MundipaggJSONAnswer: string(body)}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
