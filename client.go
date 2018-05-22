package recaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	apiEndpoint    = "https://www.google.com/recaptcha/api/siteverify"
	defaultTimeout = 15 * time.Second
)

// Option for initializer.
type Option func(*Client)

// SetHTTPClient sets httpClient.
// Used in tests for stubing.
func SetHTTPClient(httpClient *http.Client) Option {
	return func(cli *Client) {
		cli.httpClient = httpClient
	}
}

// SetTimeout sets timeout for the http client.
func SetTimeout(timeout time.Duration) Option {
	return func(cli *Client) {
		cli.httpClient.Timeout = timeout
	}
}

// Client struct to verify captcha.
type Client struct {
	secret     string
	httpClient *http.Client
}

// New returns initialized Client.
func New(secret string, options ...Option) *Client {
	cli := Client{
		secret: secret,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, option := range options {
		option(&cli)
	}

	return &cli
}

// Verify verifies reCaptcha response received from frontend.
func (cli *Client) Verify(gRecaptchaResponse string) (*Response, error) {
	return cli.verify(gRecaptchaResponse, nil)
}

// VerifyWithIP verifies reCaptcha response received from frontend with optional remoteip parameter.
func (cli *Client) VerifyWithIP(gRecaptchaResponse, remoteIP string) (*Response, error) {
	return cli.verify(gRecaptchaResponse, &remoteIP)
}

func (cli *Client) verify(gRecaptchaResponse string, remoteIP *string) (*Response, error) {
	form := url.Values{
		"secret":   []string{cli.secret},
		"response": []string{gRecaptchaResponse},
	}

	if remoteIP != nil {
		form.Set("remoteip", *remoteIP)
	}

	req, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	response := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return &response, errors.Wrap(err, "unmarshal api response")
	}

	for _, errCode := range response.ErrorCodes {
		if err, ok := errorsMap[errCode]; ok {
			return &response, err
		}
	}

	if !response.Success {
		return &response, ErrUnsucceeded
	}

	return &response, nil
}
