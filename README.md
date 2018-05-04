[![GoDoc](https://godoc.org/gopkg.in/romanyx/recaptcha.v1?status.svg)](https://godoc.org/gopkg.in/romanyx/recaptcha.v1)
[![Build Status](https://travis-ci.org/romanyx/recaptcha.png)](https://travis-ci.org/romanyx/recaptcha)
[![Go Report Card](https://goreportcard.com/badge/github.com/romanyx/recaptcha)](https://goreportcard.com/report/github.com/romanyx/recaptcha)

# recaptcha

Google's reCAPTCHA Golang implementation.

# Install

To get the package, execute:

```bash
go get gopkg.in/romanyx/recaptcha.v1
```

To import this package, add the following line to your code:

```bash
import "gopkg.in/romanyx/recaptcha.v1"
```

Refer to it as `recaptcha`.

For more details, see the API documentation.

# Example
``` go
package main

import (
	"fmt"

	import "gopkg.in/romanyx/recaptcha.v1"
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

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!
