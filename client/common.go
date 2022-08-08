package client

import (
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// linkRegex is used to extract pagination links from product search results.
var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

// Represents the pagination data from the response header
type Pagination struct {
    NextPageOptions     *ListOptions
    PreviousPageOptions *ListOptions
}

// The query string otpions provided by Shopify
type ListOptions struct {
    PageInfo     string     `url:"page_info,omitempty"`
    Page         int        `url:"page,omitempty"`
    Limit        int        `url:"limit,omitempty"`
    CreatedAtMin *time.Time `url:"created_at_min,omitempty"`
    CreatedAtMax *time.Time `url:"created_at_max,omitempty"`
    UpdatedAtMin *time.Time `url:"updated_at_min,omitempty"`
    UpdatedAtMax *time.Time `url:"updated_at_max,omitempty"`
    Order        string     `url:"order,omitempty"`
    Fields       string     `url:"fields,omitempty"`
    Vendor       string     `url:"vendor,omitempty"`
}

// extractPagination extracts pagination info from linkHeader.
// Details on the format are here: https://shopify.dev/api/usage/pagination-rest
func extractPagination(linkHeader string) (*Pagination, error) {
    pagination := new(Pagination)

    if linkHeader == "" {
        return pagination, nil
    }

    for _, link := range strings.Split(linkHeader, ",") {
        match := linkRegex.FindStringSubmatch(link)
        // Make sure the link is not empty or invalid
        if len(match) != 3 {
            // We expect 3 values:
            // match[0] = full match
            // match[1] is the URL and match[2] is either 'previous' or 'next'
            err := ResponseDecodingError{
                Message: "could not extract pagination link header",
            }
            return nil, err
        }

        rel, err := url.Parse(match[1])
        if err != nil {
            err = ResponseDecodingError{
                Message: "pagination does not contain a valid URL",
            }
            return nil, err
        }

        params, err := url.ParseQuery(rel.RawQuery)
        if err != nil {
            return nil, err
        }

        paginationListOptions := ListOptions{}

        paginationListOptions.PageInfo = params.Get("page_info")
        if paginationListOptions.PageInfo == "" {
            err = ResponseDecodingError{
                Message: "page_info is missing",
            }
            return nil, err
        }

        limit := params.Get("limit")
        if limit != "" {
            paginationListOptions.Limit, err = strconv.Atoi(params.Get("limit"))
            if err != nil {
                return nil, err
            }
        }

        // 'rel' is either next or previous
        if match[2] == "next" {
            pagination.NextPageOptions = &paginationListOptions
        } else {
            pagination.PreviousPageOptions = &paginationListOptions
        }
    }

    return pagination, nil
}
