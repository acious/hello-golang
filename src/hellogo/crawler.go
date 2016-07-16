package main

import (
	"fmt"
)


//기본이 되는 Go언어의 기본 소스 구조 작성
func main(){
	// 크롤링할 URL과 그 fetcher
    Crawl("http://golang.org/", 4, fetcher) 
}