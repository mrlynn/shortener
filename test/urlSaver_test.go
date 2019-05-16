package test

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
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
			log.Fatal(err)
		}

	}

}
