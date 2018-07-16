package recaptcha

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	secret     = "very_secret"
	okResponse = `{
		"success": true,
		"challenge_ts": "2018-07-16T08:25:45Z",
		"hostname": "http://example.com"
	}`
	unsuccessedResponse = `{
		"success": false,
		"challenge_ts": "2018-07-16T08:25:45Z",
		"hostname": "http://example.com"
	}`
	missingInputSecretResponse = `{
		"success": false,
		"error-codes": [
			"missing-input-secret"
		]
	}`
	invalidInputSecretResponse = `{
		"success": false,
		"error-codes": [
			"invalid-input-secret"
		]
	}`
)

func TestClient_verify(t *testing.T) {
	tsTime, _ := time.Parse("2006-01-02T15:04:05Z", "2018-07-16T08:25:45Z")

	type args struct {
		gRecaptchaResponse, remoteIP string
	}

	tests := []struct {
		name    string
		args    args
		handler func(args) http.Handler
		wantErr bool
		want    *Response
		err     error
	}{
		{
			name: "successed response",
			args: args{
				gRecaptchaResponse: "response_example",
				remoteIP:           "127.0.0.1",
			},
			handler: func(a args) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseForm()
					got := r.PostForm.Get("secret")
					if got != secret {
						t.Errorf("expected secret to be: %s, got: %s", secret, got)
					}
					got = r.PostForm.Get("remoteip")
					if got != a.remoteIP {
						t.Errorf("expected remoteip to be: %s, got: %s", a.remoteIP, got)
					}
					got = r.PostForm.Get("response")
					if got != a.gRecaptchaResponse {
						t.Errorf("expected response to be: %s, got: %s", a.gRecaptchaResponse, got)
					}

					fmt.Fprint(w, okResponse)
				})
			},
			wantErr: false,
			want: &Response{
				Success:     true,
				ChallengeTs: challengeTs(tsTime),
				Hostname:    "http://example.com",
			},
		},
		{
			name: "unsuccessed response",
			args: args{
				gRecaptchaResponse: "response_example",
				remoteIP:           "127.0.0.1",
			},
			handler: func(a args) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, unsuccessedResponse)
				})
			},
			wantErr: true,
			err:     ErrUnsucceeded,
			want: &Response{
				Success:     false,
				ChallengeTs: challengeTs(tsTime),
				Hostname:    "http://example.com",
			},
		},
		{
			name: "missing input secret response",
			args: args{
				gRecaptchaResponse: "response_example",
				remoteIP:           "127.0.0.1",
			},
			handler: func(a args) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, missingInputSecretResponse)
				})
			},
			wantErr: true,
			err:     ErrMissingInputSecret,
			want:    &Response{ErrorCodes: []string{"missing-input-secret"}},
		},
		{
			name: "invalid input secret response",
			args: args{
				gRecaptchaResponse: "response_example",
				remoteIP:           "127.0.0.1",
			},
			handler: func(a args) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, invalidInputSecretResponse)
				})
			},
			wantErr: true,
			err:     ErrInvalidInputSecret,
			want:    &Response{ErrorCodes: []string{"invalid-input-secret"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createFakeServer(tt.handler(tt.args))
			defer server.Close()
			httpClient := testingHTTPClient(server)

			cli := New(secret, SetHTTPClient(httpClient))
			var remoteIP *string
			if tt.args.remoteIP != "" {
				remoteIP = &tt.args.remoteIP
			}

			got, err := cli.verify(tt.args.gRecaptchaResponse, remoteIP)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if tt.wantErr && err != tt.err {
				t.Errorf("expected error: %s, got: %s", tt.err, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
