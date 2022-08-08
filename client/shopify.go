package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "path"
    "strings"
    "time"

    "github.com/google/go-querystring/query"
)

const (
    UserAgent = "goshopify/1.0.0"
)

// Client manages communication with the Shopify API.
type Client struct {
    // HTTP client used to communicate with the Shopify API.
    Client *http.Client

    // Base URL for API requests.
    // This is set on a per-store basis which means that each store must have
    // its own client.
    baseURL *url.URL

    // version you're currently using of the api, defaults to "stable"
    ApiVersion string

    // A permanent access token
    token string

    // Services used for communicating with the API
    Customer CustomerService
}

func NewClient(apiUrl, token string) *Client {
    baseURL, err := url.Parse(apiUrl)
    if err != nil {
        return nil
    }

    c := &Client{
        Client:     &http.Client{Timeout: time.Second * 60},
        baseURL:    baseURL,
        token:      token,
        ApiVersion: "2022-07",
    }

    c.Customer = &CustomerServiceOp{client: c}

    return c
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

// Post performs a POST request for the given path and saves the result in the given resource.
func (c *Client) Post(path string, data, resource interface{}) error {
    return c.CreateAndDo("POST", path, data, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the given resource.
func (c *Client) Put(path string, data, resource interface{}) error {
    return c.CreateAndDo("PUT", path, data, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(path string) error {
    return c.CreateAndDo("DELETE", path, nil, nil, nil)
}

// CreateAndDo performs a web request to Shopify with the given method (GET,
// POST, PUT, DELETE) and relative path (e.g. "/admin/orders.json").
// The data, options and resource arguments are optional and only relevant in certain situations.
// If the data argument is non-nil, it will be used as the body of the request for POST and PUT requests.
// The options argument is used for specifying request options such as search parameters like created_at_min
// Any data returned from Shopify will be marshalled into resource argument.
func (c *Client) CreateAndDo(method, relPath string, data, options, resource interface{}) error {
    _, err := c.createAndDoGetHeaders(method, relPath, data, options, resource)
    if err != nil {
        return err
    }
    return nil
}

// Creates an API request. A relative URL can be provided in urlStr, which will
// be resolved to the BaseURL of the Client. Relative URLS should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, relPath string, body, options interface{}) (*http.Request, error) {
    rel, err := url.Parse(relPath)
    if err != nil {
        return nil, err
    }

    // Make the full url based on the relative path
    u := c.baseURL.ResolveReference(rel)

    // Add custom options
    if options != nil {
        optionsQuery, err := query.Values(options)
        if err != nil {
            return nil, err
        }

        for k, values := range u.Query() {
            for _, v := range values {
                optionsQuery.Add(k, v)
            }
        }
        u.RawQuery = optionsQuery.Encode()
    }

    // A bit of JSON handling
    var js []byte = nil

    if body != nil {
        js, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }

    req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
    if err != nil {
        return nil, err
    }

    req.Header.Add("User-Agent", UserAgent)
    req.Header.Add("X-Shopify-Access-Token", c.token)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/json")

    return req, nil
}

// createAndDoGetHeaders creates an executes a request while returning the response headers.
func (c *Client) createAndDoGetHeaders(method, relPath string, data, options, resource interface{}) (http.Header, error) {
    if strings.HasPrefix(relPath, "/") {
        // make sure it's a relative path
        relPath = strings.TrimLeft(relPath, "/")
    }

    prefix := fmt.Sprintf("%s/%s", "api", c.ApiVersion)

    relPath = path.Join(prefix, relPath)
    req, err := c.NewRequest(method, relPath, data, options)
    if err != nil {
        return nil, err
    }

    return c.doGetHeaders(req, resource)
}

// doGetHeaders executes a request, decoding the response into `resp` and also returns any response headers.
func (c *Client) doGetHeaders(req *http.Request, resp interface{}) (http.Header, error) {
    result, err := c.Client.Do(req)
    if err != nil {
        return nil, err
    }

    defer result.Body.Close()

    respErr := CheckResponseError(result)
    if respErr != nil {
        return nil, respErr
    }

    if resp == nil {
        return nil, nil
    }

    decoder := json.NewDecoder(result.Body)
    err = decoder.Decode(&resp)
    if err != nil {
        return nil, err
    }

    return result.Header, nil
}
