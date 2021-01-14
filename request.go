package mundipagg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lusantisuper/mundipagg/internal/utils"
)

// MakePostRequest do the low-level request
func MakePostRequest(data interface{}, secretKey string, indepotencyKey string, url string) (*Response, error) {
	postData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))

	// Setting the headers to make the request
	req.Header.Set("Authorization", "Basic "+utils.ToBase64(secretKey+":"))
	req.Header.Set("Content-Type", "application/json")
	if !utils.IsStringEmpty(indepotencyKey) {
		req.Header.Set("Idempotency-Key", indepotencyKey)
	}

	// Running the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Results of the request
	body, _ := ioutil.ReadAll(resp.Body)

	respCode, err := strconv.Atoi(resp.Status)
	if err != nil {
		return nil, errors.New("Wrong response code")
	}

	if respCode < 200 || respCode >= 400 {
		return nil, errors.New("Invalid Request:\nSended:\n" + string(postData) + "Received:\n" + string(body))
	}

	// Saving the result
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}