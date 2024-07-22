package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("https://pkg.go.dev/search?m=package&q=http")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	parse(body)
}

func walkBlock(tkn *html.Tokenizer, callback func(html.TokenType, *html.Tokenizer)) {
	depth := 1

	for {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			return

		case html.EndTagToken:
			depth--
			if depth == 0 {
				return
			}

		case html.StartTagToken:
			depth++
		}

		callback(tt, tkn)
	}
}

func parseTitle(tkn *html.Tokenizer) {
	title := ""
	link := ""

	walkBlock(tkn, func(tt html.TokenType, tkn *html.Tokenizer) {
		switch tt {
		case html.StartTagToken:
			t := tkn.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						link = a.Val
					}
				}
				tt = tkn.Next()
				t = tkn.Token()
				if tt == html.TextToken {
					title = strings.Trim(t.Data, " \n\t")
				}
			}
		}
	})

	fmt.Println(title, link)
}

func parseSnippet(tkn *html.Tokenizer) {
	walkBlock(tkn, func(tt html.TokenType, tkn *html.Tokenizer) {
		switch tt {
		case html.StartTagToken:
			t := tkn.Token()
			if t.Data == "h2" {
				parseTitle(tkn)
			}
		}
	})
}

func parse(text []byte) {
	tkn := html.NewTokenizer(bytes.NewReader(text))
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := tkn.Token()
			for _, a := range t.Attr {
				if a.Key == "class" && a.Val == "SearchSnippet" {
					parseSnippet(tkn)
					return
				}
			}

		}
	}
}
