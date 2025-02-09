package pageutils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
)

const (
	linkFormat = "%v%v?%v"
)

func NewLinks(r *http.Request, page, size, totalItem, totalPage int) *dto.Links {
	queries := r.URL.Query()
	host := r.Host
	path := r.URL.Path
	if r.TLS != nil {
		host = fmt.Sprintf("https://%v", host)
	} else {
		host = fmt.Sprintf("http://%v", host)
	}

	queries.Set("page", strconv.Itoa(page))
	selfLink := fmt.Sprintf(linkFormat, host, path, queries.Encode())

	queries.Set("page", "1")
	firstLink := fmt.Sprintf(linkFormat, host, path, queries.Encode())

	if totalPage > 0 {
		queries.Set("page", strconv.Itoa(totalPage))
	} else {
		queries.Set("page", "1")
	}
	lastLink := fmt.Sprintf(linkFormat, host, path, queries.Encode())

	return &dto.Links{
		Self:  selfLink,
		First: firstLink,
		Prev:  createPrevLink(queries, host, path, page),
		Next:  createNextLink(queries, host, path, page, totalPage),
		Last:  lastLink,
	}
}

func createNextLink(queries url.Values, host, path string, page, totalPage int) string {
	if page+1 >= totalPage {
		if totalPage > 0 {
			queries.Set("page", strconv.Itoa(totalPage))
		} else {
			queries.Set("page", "1")
		}
	} else {
		queries.Set("page", strconv.Itoa(page+1))
	}

	return fmt.Sprintf(linkFormat, host, path, queries.Encode())
}

func createPrevLink(queries url.Values, host, path string, page int) string {
	if page-1 <= 0 {
		queries.Set("page", "1")
	} else {
		queries.Set("page", strconv.Itoa(page-1))
	}

	return fmt.Sprintf(linkFormat, host, path, queries.Encode())
}
