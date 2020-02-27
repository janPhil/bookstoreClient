package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

func NewClient(baseURL *url.URL, client *http.Client) *Client {
	return &Client{
		BaseURL:    baseURL,
		httpClient: client,
	}
}

func (c *Client) ListBooks() (Books, error) {
	rel := &url.URL{Path: "/books"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var books Books
	err = json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (c *Client) CreateBook(b *Book) error {
	requestBody, err := json.Marshal(b)
	if err != nil {
		return err
	}

	rel := &url.URL{Path: "/book"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}

	if res.StatusCode == http.StatusNotFound {
		return errors.New("Book not found")
	}
	return nil
}

func (c *Client) ListBook(isbn string) (Book, error) {
	rel := &url.URL{Path: "/books/" + isbn}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Book{}, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return Book{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var book Book
		err = json.NewDecoder(res.Body).Decode(&book)
		if err != nil {
			return Book{}, err
		}
		return book, nil
	}

	if res.StatusCode == http.StatusNotFound {
		return Book{}, errors.New("Book not found")
	}
	return Book{}, err

}
