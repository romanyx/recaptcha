[![GoDoc](https://godoc.org/github.com/romanyx/recaptcha?status.svg)](https://godoc.org/github.com/romanyx/recaptcha)
[![Build Status](https://travis-ci.org/romanyx/recaptcha.png)](https://travis-ci.org/romanyx/recaptcha)

# recaptcha

Google's reCAPTCHA Golang implementation.

# Usage
``` go
package main

import (
	"fmt"

	"github.com/romanyx/recaptcha"
)

func main() {
	r := recaptcha.New("secret")
	res, err := r.Verify("response") // g-recaptcha-response parameter
	if err != nil {
		switch err {
		case recaptcha.ErrMissingInputSecret:
			fmt.Println(err)
		case recaptcha.ErrInvalidInputSecret:
			fmt.Println(err)
		case recaptcha.ErrMissingInputResponse:
			fmt.Println(err)
		case recaptcha.ErrInvalidInputResponse:
			fmt.Println(err)
		case recaptcha.ErrBadRequest:
			fmt.Println(err)
		case recaptcha.ErrUnsucceeded:
			// This triggers when response.Success is equal false.
			fmt.Println(err)
		default:
			fmt.Printf("unknown error: %s\n", err)
		}
	}

	// If err is equal to nil, then verification has been successed.
	if err == nil {
		fmt.Printf("%+v\n", res)
	}
}
```