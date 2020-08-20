package main

import (
	"github.com/go-resty/resty"
	"testing"
)

//testing status 200 using resty
func Test_GetUs90210_StatusCodeShouldEqual200(t *testing.T) {
	client := resty.New()
	resp, _ := client.R().Get("http://localhost:8000/persons")
	if resp.StatusCode() != 200 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 200, resp.StatusCode())
	}
}
