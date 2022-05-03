package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"

	"strings"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
			break
		}
	}
	return
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)
	defer func() {
		chFinished <- true
	}()
	if err != nil {
		return
	}
	body := resp.Body
	defer body.Close()
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				ok, val := getHref(t)
				if !ok {
					continue
				}
				if strings.Index(val, "http") == 0 {
					ch <- val
				}

			}
		}

	}

}

func main() {
	ch := make(chan string)
	chFinished := make(chan bool)
	seedURLs := []string{"https://insomnia.rest", "https://twitter.com/GregorySchier", "https://support.insomnia.rest", "https://chat.insomnia.rest"}
	foundURLs := make(map[string]bool)

	for _, url := range seedURLs {
		go crawl(url, ch, chFinished)
	}

	for i := 0; i < len(seedURLs); {
		select {
		case url := <-ch:
			foundURLs[url] = true
		case <-chFinished:
			i++

		}

	}

	fmt.Println(foundURLs)

}
