package main

import (
	"fmt"
	"sync"
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

	visited.Lock()	//lock the map
	if visited.urls[url] == true {	//if this url is already visited, return
		visited.Unlock()
		return
	}
	visited.urls[url] = true
	visited.Unlock()	//release the lock on map

	body, urls, err := fetcher.Fetch(url)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan bool)		//create a channel to wait for the completion of all go routines

	for _, nestedUrl := range urls {
		go func(ur string) {
			Crawl(ur, depth-1, fetcher)
			done <- true
		}(nestedUrl)
	}


	for i := range urls {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls))
		<-done
	}
	fmt.Printf("<- Done with %v\n", url)
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
    

	fmt.Println("Fetching stats\n--------------")

	for url, err := range visited.urls {
		if err != true {
			fmt.Printf("%v failed: %v\n", url, err)
		} else {
			fmt.Printf("%v was fetched\n", url)
		}
	}
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
