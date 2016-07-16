package main

import (
	"fmt"
)

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}
