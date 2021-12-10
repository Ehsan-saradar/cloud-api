package zarinpal

import (
	"bytes"
	"encoding/json"
	"net/http"

	"api.cloud.io/pkg/payment"
	"github.com/pkg/errors"
)

const (
	mainAPIEndpoint            = "https://www.zarinpal.com/pg/rest/WebGate/"
	mainPayEndpoint            = "https://www.zarinpal.com/pg/StartPay/"
	testAPIEndpoint            = "https://sandbox.zarinpal.com/pg/rest/WebGate/"
	testPayEndpoint            = "https://sandbox.zarinpal.com/pg/StartPay/"
	paymentRequestCommand      = "PaymentRequest.json"
	paymentVerificationCommand = "PaymentVerification.json"
)

var _ payment.Provider = (*Provider)(nil)

// Provider implements payment.Provider interface for Zarnpal webpay.
// It supports "email" and "mobile" meta for Request.
type Provider struct {
	merchId     string
	apiEndpoint string
	payEndpoint string
	sandbox     bool
	httpCli     *http.Client
}

// NewProvider returns a new Provider from the given merchant id for main and sandbox platform.
func NewProvider(merchId string, sandbox bool) (*Provider, error) {
	if len(merchId) != 36 {
		return nil, errors.New("merchId must be 36 characters")
	}

	apiEndpoint := mainAPIEndpoint
	payEndpoint := mainPayEndpoint
	if sandbox == true {
		apiEndpoint = testAPIEndpoint
		payEndpoint = testPayEndpoint
	}

	return &Provider{
		merchId:     merchId,
		apiEndpoint: apiEndpoint,
		payEndpoint: payEndpoint,
		sandbox:     sandbox,
		httpCli:     &http.Client{},
	}, nil
}

type paymentRequestReq struct {
	MerchantId  string `json:"MerchantID"`
	Amount      int64  `json:"Amount"`
	CallbackURL string `json:"CallbackURL"`
	Description string `json:"Description"`
	Email       string `json:"Email"`
	Mobile      string `json:"Mobile"`
}

type paymentRequestRes struct {
	Status    int64  `json:"Status"`
	Authority string `json:"Authority"`
}

// New is implementation of Provider.New
func (zp *Provider) New(req *payment.Request) (*payment.Session, error) {
	err := zp.validateRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "payment request is not valid for zarinpal service")
	}

	payload := paymentRequestReq{
		MerchantId:  zp.merchId,
		Amount:      req.Amount,
		CallbackURL: req.CallbackURL,
		Description: req.Description,
		Email:       getMetaField(req.Meta, "email"),
		Mobile:      getMetaField(req.Meta, "mobile"),
	}
	var res paymentRequestRes
	err = zp.request(paymentRequestCommand, &payload, &res)
	if err != nil {
		return nil, err
	}

	if res.Status == 100 {
		s := &payment.Session{
			Key:        res.Authority,
			GatewayURL: zp.payEndpoint + res.Authority,
			Amount:     req.Amount,
		}
		return s, nil
	}

	return nil, errors.Errorf("zarinpal returned a non-seccessfull response with status code %d", res.Status)
}

type paymentVerificationReq struct {
	MerchantId string `json:"MerchantID"`
	Authority  string `json:"Authority"`
	Amount     int64  `json:"Amount"`
}

type paymentVerificationRes struct {
	Status int64       `json:"Status"`
	RefId  json.Number `json:"RefID"`
}

// Verify implements Provider.Verify
func (zp *Provider) Verify(s *payment.Session) (string, error) {
	err := zp.validateSession(s)
	if err != nil {
		return "", errors.Wrap(err, "payment session is not valid for zarinpal service")
	}

	payload := paymentVerificationReq{
		MerchantId: zp.merchId,
		Authority:  s.Key,
		Amount:     s.Amount,
	}
	var res paymentVerificationRes
	err = zp.request(paymentVerificationCommand, &payload, &res)
	if err != nil {
		return "", err
	}

	if res.Status == 100 {
		return res.RefId.String(), nil
	}

	return "", errors.Errorf("zarinpal returned a non-seccessfull response with status code %d", res.Status)
}

func (zp *Provider) validateRequest(req *payment.Request) error {
	if req.Amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if req.CallbackURL == "" {
		return errors.New("callbackURL should not be empty")
	}
	if req.Description == "" {
		return errors.New("description should not be empty")
	}
	return nil
}

func (zp *Provider) validateSession(s *payment.Session) error {
	if s.Amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if s.Key == "" {
		return errors.New("key should not be empty")
	}
	return nil
}

func (zp *Provider) request(command string, payload interface{}, res interface{}) error {
	url := zp.apiEndpoint + command

	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return errors.Wrap(err, "could not encode payload to json")
	}

	resp, err := zp.httpCli.Post(url, "application/json", body)
	if err != nil {
		return errors.Wrap(err, "http request to zarinpal servers failed")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return errors.Wrap(err, "could not decode the response to json")
	}
	return nil
}

func getMetaField(m map[string]string, key string) string {
	if m == nil {
		return ""
	}

	if str, ok := m[key]; ok {
		return str
	}
	return ""
}
