package client

import (
    "net/http"
    "net/url"
)

// Client manages communication with the Shopify API.
type Client struct {
    // HTTP client used to communicate with the Shopify API.
    Client *http.Client

    // Base URL for API requests.
    // This is set on a per-store basis which means that each store must have
    // its own client.
    baseURL *url.URL

    // URL Prefix, defaults to "admin" see WithVersion
    pathPrefix string

    // version you're currently using of the api, defaults to "stable"
    ApiVersion string

    // A permanent access token
    token string

    // Services used for communicating with the API
    Customer CustomerService
}

func (c *Client) Count(path string, options interface{}) (int, error) {
    resource := struct {
        Count int `json:"count"`
    }{}
    err := c.Get(path, &resource, options)
    return resource.Count, err
}

// Get performs a GET request for the given path and saves the result in the given resource.
func (c *Client) Get(path string, resource, options interface{}) error {
    return c.CreateAndDo("GET", path, nil, options, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *Client) Post(path string, data, resource interface{}) error {
    return c.CreateAndDo("POST", path, data, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *Client) Put(path string, data, resource interface{}) error {
    return c.CreateAndDo("PUT", path, data, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(path string) error {
    return c.CreateAndDo("DELETE", path, nil, nil, nil)
}

// CreateAndDo performs a web request to Shopify with the given method (GET,
// POST, PUT, DELETE) and relative path (e.g. "/admin/orders.json").
// The data, options and resource arguments are optional and only relevant in
// certain situations.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
// The options argument is used for specifying request options such as search
// parameters like created_at_min
// Any data returned from Shopify will be marshalled into resource argument.
func (c *Client) CreateAndDo(method, relPath string, data, options, resource interface{}) error {
    _, err := c.createAndDoGetHeaders(method, relPath, data, options, resource)
    if err != nil {
        return err
    }
    return nil
}
