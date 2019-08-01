// Package test provides testing coverage for the web interface of this
// example package. This package depends upon, and uses the net/http package
// to facilitate a web-based interaction with the shortener to test its
// functionality.
package test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	fuzz "github.com/google/gofuzz"
)

// TestSaving is a function that drives a web interaction to test the shortener. It
// uses a text file "urls.txt" that contains a list of urls to be incorporated in the test.
func TestSaving(t *testing.T) {
	file, err := os.Open("urls.txt")

	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		_, err := http.PostForm("http://localhost:8080/shortener", url.Values{
			"url": {scanner.Text()},
		})

		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestValidation is a function that tests the correctness of the shortener.
func TestValidation(t *testing.T) {
	f := fuzz.New()

	var testUrl string

	f.Fuzz(&testUrl)

	if !strings.Contains(testUrl, "http") {
		testUrl = "http://" + testUrl
	}

	resp, err := http.PostForm("http://localhost:8080/shortener", url.Values{
		"url": {testUrl},
	})

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(testUrl, string(body))
}
