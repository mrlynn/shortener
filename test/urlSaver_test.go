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
