# requestworker


[![Go Report Card](https://goreportcard.com/badge/github.com/z9905080/requestworker)](https://goreportcard.com/report/github.com/z9905080/requestworker)
[![Build Status](https://travis-ci.org/z9905080/requestworker.svg?branch=master)](https://travis-ci.org/z9905080/requestworker)

a lib for go to batch processing send web request

## Install

`go get github.com/z9905080/requestworker`

### Usage

```

func main() {

    // Init request
	req, err := http.NewRequest("GET", "http://tw.yahoo.com", nil)
	if err != nil {
		t.Error("request error: ", err)
	}

	// Init worker
	a := New(5)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = a.Execute(ctx, req, func(resp *http.Response, err error) error {

		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil

	})
}

```
