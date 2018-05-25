package recaptcha

import (
	"errors"
	"strings"
	"time"
)

var (
	// ErrMissingInputSecret when the secret parameter is missing.
	ErrMissingInputSecret = errors.New("missing-input-secret")
	// ErrInvalidInputSecret when the secret parameter is invalid or malformed.
	ErrInvalidInputSecret = errors.New("invalid-input-secret")
	// ErrMissingInputResponse when the response parameter is missing.
	ErrMissingInputResponse = errors.New("missing-input-response")
	// ErrInvalidInputResponse when the response parameter is invalid or malformed.
	ErrInvalidInputResponse = errors.New("invalid-input-response")
	// ErrBadRequest when the request is invalid or malformed.
	ErrBadRequest = errors.New("bad-request")
	// ErrUnsucceeded when response success value equals to false.
	ErrUnsucceeded = errors.New("unsucceeded status")

	errorsMap = map[string]error{
		"missing-input-secret":   ErrMissingInputSecret,
		"invalid-input-secret":   ErrInvalidInputSecret,
		"missing-input-response": ErrMissingInputResponse,
		"invalid-input-response": ErrInvalidInputResponse,
		"bad-request":            ErrBadRequest,
	}
)

// Response holds data from reCaptcha API response.
type Response struct {
	Success     bool        `json:"success"`
	ChallengeTs challengeTs `json:"challenge_ts"`
	Hostname    string      `json:"hostname"`
	ErrorCodes  []string    `json:"error-codes"`
}

type challengeTs time.Time

func (t *challengeTs) UnmarshalJSON(data []byte) error {
	asString := strings.Trim(string(data), `"`)
	parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", asString)
	if err != nil {
		return err
	}
	*t = challengeTs(parsedTime)
	return err
}

func (t *challengeTs) String() string {
	return time.Time(*t).String()
}
