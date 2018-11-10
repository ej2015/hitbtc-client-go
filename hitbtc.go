package hitbtc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HTTPClient *http.Client
}

type Symbol struct {
	Id                   string
	BaseCurrency         string
	QuoteCurrency        string
	QuantityIncrement    float64
	TickSize             float64
	TakeLiquidityRate    float64
	ProvideLiquidityRate float64
	FeeCurrency          string
}

type Currency struct {
	Id                 string
	FullName           string `json:"fullName"`
	Crypto             bool   `json:"crypto"`
	PayinEnabled       bool   `json:"payinEnabled"`
	PayinPaymentId     bool   `json:"payinPaymentId"`
	PayinConfirmations int64  `json:"payinConfirmations"`
	PayoutEnabled      bool   `json:"payoutEnabled"`
	PayoutIsPaymentId  bool   `json:"payoutIsPaymentId"`
	TransferEnabled    bool   `json:"transferEnabled"`
	Delisted           bool   `json:"delisted"`
	PayoutFee          string `json:"payoutFee"`
}

func (c *Client) Symbols() (interface{}, error) {
	req, err := c.newRequest("GET", "public/symbol", nil)

	if err != nil {
		return nil, err
	}
	var res []Symbol
	resp, err := c.do(req, res)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *Client) Symbol(symbol string) (interface{}, error) {
	path := path.Join("public/symbol", symbol)
	req, err := c.newRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}
	var res Symbol
	resp, err := c.do(req, res)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *Client) Currencies() (interface{}, error) {
	req, err := c.newRequest("GET", "public/currency", nil)

	if err != nil {
		return nil, err
	}

	var res []Currency
	///var resp *http.Response
	resp, err := c.do(req, res)
	//err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	return resp, err

}

func (c *Client) Currency(currency string) (interface{}, error) {
	path := path.Join("public/currency", currency)
	req, err := c.newRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	var res Currency
	resp, err := c.do(req, res)

	if err != nil {
		return nil, err
	}

	return resp, err

}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

//func (c *Client) do(req *http.Request) (*http.Response, error) {
func (c *Client) do(req *http.Request, v interface{}) (interface{}, error) {
	resp, err := c.HTTPClient.Do(req)
	log.Println(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
