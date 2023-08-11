package aoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrorInfo is a struct of error information.
type ErrorInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// ResponseError is a struct of response error.
type ResponseError struct {
	E ErrorInfo `json:"error"`
}

// Error returns the error message.
func (respErr *ResponseError) Error() string {
	return fmt.Sprintf("type=%q, param=%q, code=%q: %s", respErr.E.Type, respErr.E.Param, respErr.E.Code, respErr.E.Message)
}

// build builds the error from the response. It always returns an error.
func (respErr *ResponseError) build(reader io.ReadCloser, statusCode int) error {
	defer func() {
		_ = reader.Close()
	}()

	err := errors.Join(ErrResponse, fmt.Errorf("status code %d", statusCode))

	if e := json.NewDecoder(reader).Decode(respErr); e != nil {
		return errors.Join(err, fmt.Errorf("failed unmarshal error: %w", e))
	}

	return errors.Join(err, respErr)
}

// Params is a struct of API authentication and additional parameters.
type Params struct {
	Bearer       string
	Organization string
	URL          string
	StopMarker   string
}

// CommonRequest is a common interface for all API requests.
type CommonRequest interface {
	build(ctx context.Context, auth *Params) (*http.Request, error)
}

// commonRequest sends a request to the API and returns a body response.
// A caller must close the response body.
func commonRequest(ctx context.Context, client *http.Client, cReq CommonRequest, p Params) (io.ReadCloser, error) {
	request, err := cReq.build(ctx, &p)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respErr := &ResponseError{}
		// build closes the response body
		return nil, respErr.build(resp.Body, resp.StatusCode)
	}

	return resp.Body, nil
}
