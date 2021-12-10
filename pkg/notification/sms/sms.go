package sms

import (
	"fmt"
	"github.com/kavenegar/kavenegar-go"
)
var Cli *Client
type Client struct {
	api *kavenegar.Kavenegar
	keyTemplate string
}
func NewClient(apiKey  ,kycTemplate string) *Client{
	api := kavenegar.New(apiKey)
	Cli:=&Client{
		api,
		kycTemplate,
	}
	return Cli
}

func (cli *Client) SendOtp(receptor string,token string) error {
	fmt.Println(receptor+":"+token)
	params := &kavenegar.VerifyLookupParam{
	}
	if _, err := cli.api.Verify.Lookup(receptor, cli.keyTemplate, token, params); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			return err
		case *kavenegar.HTTPError:
			return err
		default:
			return err
		}
	}
	return nil
}