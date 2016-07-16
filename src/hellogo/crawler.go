package main

import (
	"fmt"
)

type fakeResult struct {
	body string
	urls []string
}

type fakeFetcher map[string]*fakeResult

var fetcher = &fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

type result struct {
	url, body string
	urls      []string
	err       error
	depth     int
}

func Crawl(url string, depth int, fetcher Fetcher) {
	results := make(chan *result)
	fetched := make(map[string]bool)
	fetch := func(url string, depth int) {
		body, urls, err := fetcher.Fetch(url)
		results <- &result{url, body, urls, err, depth}
	}

	go fetch(url, depth)
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		res := <-results

		if res.err != nil {
			fmt.Println(res.err)
			continue
		}

		fmt.Printf("found : %s %q\n", res.url, res.body)

		if res.depth > 0 {
			for _, url := range res.urls {
				if !fetched[url] {
					fetching++
					go fetch(url, res.depth-1)
					fetched[url] = true
				}
			}
		}
	}

	close(results)
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}
