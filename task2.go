package main

import (
	"fmt"
	"sync"
//	"time"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

var visited = struct {
	urls map[string]bool
	sync.Mutex
}{urls: make(map[string]bool)}

func Crawl(url string, depth int, fetcher Fetcher) {

	if depth <= 0 {
		return
	}

	visited.Lock()
	if visited.urls[url] == true {
		visited.Unlock()
		return
	}
	visited.urls[url] = true
	visited.Unlock()

	body, urls, err := fetcher.Fetch(url)

	fmt.Println("main body: ", body)
	fmt.Println("main urls: ", urls)
	fmt.Println("main error: ", err)
	
	visited.Lock()
 	visited.urls[url] = true
        visited.Unlock()
	
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan bool)

	for _, nestedUrl := range urls {
		go func(ur string) {
			fmt.Printf("-> Crawling child %v of %v with depth %v \n", nestedUrl, url, depth)
			Crawl(ur, depth-1, fetcher)
			done <- true
		}(nestedUrl)
	}


	for i := range urls {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls))
		<-done
	}
	fmt.Printf("<- Done with %v\n", url)
	//time.Sleep(time.Second * 10)
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
    
//	time.Sleep(time.Second * 10)

	fmt.Println("Fetching stats\n--------------")

	for url, err := range visited.urls {
		if err != true {
			fmt.Printf("%v failed: %v\n", url, err)
		} else {
			fmt.Printf("%v was fetched\n", url)
		}
	}
//	fmt.Printf("visited ", visited.urls)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls     []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
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
