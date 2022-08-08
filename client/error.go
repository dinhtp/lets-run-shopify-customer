package client

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "reflect"
    "sort"
    "strconv"
    "strings"
)

// A general response error that follows a similar layout to Shopify response
// errors, i.e. either a single message or a list of messages.
type ResponseError struct {
    Status  int
    Message string
    Errors  []string
}

// An error specific to a rate-limiting response. Embeds the ResponseError to
// allow consumers to handle it the same was a normal ResponseError.
type RateLimitError struct {
    ResponseError
    RetryAfter int
}

// ResponseDecodingError occurs when the response body from Shopify could not be parsed.
type ResponseDecodingError struct {
    Body    []byte
    Message string
    Status  int
}

// GetStatus returns http  response status
func (e ResponseError) GetStatus() int {
    return e.Status
}

// GetMessage returns response error message
func (e ResponseError) GetMessage() string {
    return e.Message
}

// GetErrors returns response errors list
func (e ResponseError) GetErrors() []string {
    return e.Errors
}

func (e ResponseError) Error() string {
    if e.Message != "" {
        return e.Message
    }

    sort.Strings(e.Errors)
    s := strings.Join(e.Errors, ", ")

    if s != "" {
        return s
    }

    return "Unknown Error"
}

func (e ResponseDecodingError) Error() string {
    return e.Message
}

func CheckResponseError(r *http.Response) error {
    if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
        return nil
    }

    // Create an anonymously struct to parse the JSON data into.
    shopifyError := struct {
        Error  string      `json:"error"`
        Errors interface{} `json:"errors"`
    }{}

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return err
    }

    // empty body, this probably means shopify returned an error with no body
    // we'll handle that error in wrapSpecificError()
    if len(bodyBytes) > 0 {
        err := json.Unmarshal(bodyBytes, &shopifyError)
        if err != nil {
            fmt.Println("Response body:", string(bodyBytes))
            return ResponseDecodingError{
                Body:    bodyBytes,
                Message: err.Error(),
                Status:  r.StatusCode,
            }
        }
    }

    // Create the response error from the Shopify error.
    responseError := ResponseError{
        Status:  r.StatusCode,
        Message: shopifyError.Error,
    }

    // If the errors field is not filled out, we can return here.
    if shopifyError.Errors == nil {
        return wrapSpecificError(r, responseError)
    }

    // Shopify errors usually have the form:
    // {
    //   "errors": {
    //     "title": [
    //       "something is wrong"
    //     ]
    //   }
    // }
    // This structure is flattened to a single array:
    // [ "title: something is wrong" ]
    //
    // Unfortunately, "errors" can also be a single string so we have to deal
    // with that. Lots of reflection :-(
    switch reflect.TypeOf(shopifyError.Errors).Kind() {
    case reflect.String:
        // Single string, use as message
        responseError.Message = shopifyError.Errors.(string)
    case reflect.Slice:
        // An array, parse each entry as a string and join them on the message
        // json always serializes JSON arrays into []interface{}
        for _, elem := range shopifyError.Errors.([]interface{}) {
            responseError.Errors = append(responseError.Errors, fmt.Sprint(elem))
        }
        responseError.Message = strings.Join(responseError.Errors, ", ")
    case reflect.Map:
        // A map, parse each error for each key in the map.
        // json always serializes into map[string]interface{} for objects
        for k, v := range shopifyError.Errors.(map[string]interface{}) {
            // Check to make sure the interface is a slice
            // json always serializes JSON arrays into []interface{}
            if reflect.TypeOf(v).Kind() == reflect.Slice {
                for _, elem := range v.([]interface{}) {
                    // If the primary message of the response error is not set, use
                    // any message.
                    if responseError.Message == "" {
                        responseError.Message = fmt.Sprintf("%v: %v", k, elem)
                    }
                    topicAndElem := fmt.Sprintf("%v: %v", k, elem)
                    responseError.Errors = append(responseError.Errors, topicAndElem)
                }
            }
        }
    }

    return wrapSpecificError(r, responseError)
}

// Shopify reference: https://shopify.dev/api/usage/response-codes
func wrapSpecificError(r *http.Response, err ResponseError) error {
    if err.Status == http.StatusTooManyRequests {
        f, _ := strconv.ParseFloat(r.Header.Get("Retry-After"), 64)
        return RateLimitError{
            ResponseError: err,
            RetryAfter:    int(f),
        }
    }

    if err.Status == http.StatusNotAcceptable {
        err.Message = http.StatusText(err.Status)
    }

    return err
}
