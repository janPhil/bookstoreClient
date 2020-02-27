package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/janPhil/bookstoreClient/internal"
)

func main() {

	httpClient := &http.Client{Timeout: time.Duration(5 * time.Second)}
	baseURL, err := url.Parse("http://localhost:9090")
	if err != nil {
		log.Fatal("unable to parse url")
	}

	booksClient := internal.NewClient(baseURL, httpClient)

	b, err := booksClient.ListBooks()
	if err != nil {
		log.Fatal("Couldnt receive any books", err)
	}

	for _, book := range b {
		fmt.Println(book.Title)
	}

	book := &internal.Book{Author: "Jan-Philipp Heinrich", Isbn: "123-456-789", Price: 19.99, Title: "Hello World"}

	err = booksClient.CreateBook(book)
	if err != nil {
		log.Fatal("Failed to create book ", err)
	}

	bk, err := booksClient.ListBook("123-456-789")
	if err != nil {
		log.Fatal("Failed to get Book: ", err)
	}

	fmt.Println(bk)

}
