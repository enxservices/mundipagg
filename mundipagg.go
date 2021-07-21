package mundipagg

import (
	"net/http"
	"time"
)

type Mundipagg interface {
	NewSubscription(subscription *Subscription, uuid string) (*Response, error)
	NewCustomer(customer *Customer, uuid string) (*Response, error)
	NewCardByToken(customerID string, cardToken string, uuid string) (*Response, error)
	UpdateStartAt(input *string, subscriptionID string, uuid string) (*Response, error)
	UpdateNextBillingDay(nextBillingDay *time.Time, customerID string, uuid string) (*Response, error)
	AddDiscount(billExtras *BillExtras, subscriptionID string, uuid string) (*Response, error)
}

type mundipagg struct {
	BasicSecretAuthKey string
}

func New(secret string) Mundipagg {
	return mundipagg{
		BasicSecretAuthKey: secret,
	}
}

func (m mundipagg) NewCardByToken(customerID string, cardToken string, uuid string) (*Response, error) {
	creditCardLink := BASEURL + customerID + CARDENDPOINT
	card := &CreditCardToken{
		Token: cardToken,
	}

	resp, err := Do(http.MethodPost, card, m.BasicSecretAuthKey, uuid, creditCardLink)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m mundipagg) NewCustomer(c *Customer, uuid string) (*Response, error) {
	resp, err := Do(http.MethodPost, c, m.BasicSecretAuthKey, uuid, CUSTOMERURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m mundipagg) AddDiscount(extras *BillExtras, subscriptionID string, uuid string) (*Response, error) {
	completeURL := SUBSCRIPTIONURL + "/" + subscriptionID + DISCOUNTENDPOINT
	resp, err := Do(http.MethodPost, extras, m.BasicSecretAuthKey, uuid, completeURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m mundipagg) NewSubscription(s *Subscription, uuid string) (*Response, error) {
	resp, err := Do(http.MethodPost, s, m.BasicSecretAuthKey, uuid, SUBSCRIPTIONURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m mundipagg) UpdateStartAt(startAt *string, subscriptionID string, uuid string) (*Response, error) {
	input := struct {
		StartAt *string `json:"start_at,omitempty"`
	}{
		StartAt: startAt,
	}
	completeURL := SUBSCRIPTIONURL + "/" + subscriptionID + SUBSCRIPTIONUPDATESTARTATURL
	resp, err := Do(http.MethodPatch, input, m.BasicSecretAuthKey, uuid, completeURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m mundipagg) UpdateNextBillingDay(nextBillingDay *time.Time, customerID string, indepotencyKey string) (*Response, error) {
	input := struct {
		NextBillingDay *time.Time `json:"next_billing_at,omitempty"`
	}{
		NextBillingDay: nextBillingDay,
	}
	completeURL := SUBSCRIPTIONURL + "/" + customerID + SUBSCRIPTIONUPDATENEXTBILLINGDAYURL
	resp, err := Do(http.MethodPatch, input, m.BasicSecretAuthKey, indepotencyKey, completeURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
