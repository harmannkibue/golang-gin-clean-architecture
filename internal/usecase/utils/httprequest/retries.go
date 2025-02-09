package httprequest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	originalBody []byte
	err          error
)

// RetryDo when given a request it executes the request client with the specified redoes -.
func (h HttpRequest) RetryDo(req *http.Request, maxRetries int, timeout time.Duration,
	backoffStrategy []time.Duration) (*http.Response, error) {

	if req != nil && req.Body != nil {
		originalBody, err = h.copyBody(req.Body)
		h.resetBody(req, originalBody)
	}

	if err != nil {
		return nil, err
	}

	if maxRetries <= 0 {
		maxRetries = 1
	}

	client := http.Client{
		Timeout: timeout,
	}

	var resp *http.Response

	// loop number of retries -.
	for i := 1; i <= maxRetries; i++ {

		resp, err = client.Do(req)

		if err != nil {

			// Try reinitializing the request client in case of a failure -.
			if i < maxRetries {
				continue
			}
			log.Printf("failed initializing http request client. %s \n", err.Error())
			return nil, fmt.Errorf("failed initialising request %s ", err.Error())
		}

		resp, err := h.analyseRespErrors(req, resp)

		if err != nil {
			return nil, err
		}

		// The request is successful and exit the for loop with a return -.
		if resp != nil && err == nil {
			return resp, nil
		}
		// resetting body -.
		h.resetBody(req, originalBody)

		time.Sleep(backoffStrategy[i-1] + 1*time.Microsecond)
	}

	// Here, it means that retrying is useless -.
	return nil, fmt.Errorf("something went wrong please try again later")
}

// analyseRespErrors check the response status codes and chose on which to retry with -.
func (h HttpRequest) analyseRespErrors(req *http.Request, resp *http.Response) (*http.Response, error) {

	var processedResponse *http.Response
	var err error
	// The status code is withing the 400-499 range and thus contains error message -.
	if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
		processedResponse = nil
		err = fmt.Errorf("\nfailed with status code %d \n", resp.StatusCode)
	} else if resp.StatusCode < 400 {
		log.Printf("\nSUCCESSFULLY SENT REQUEST TO %s \n", req.URL)
		processedResponse = resp
		err = nil

	} else if resp.StatusCode >= 500 {
		processedResponse = nil
		err = nil
	}

	return processedResponse, err
}

// copyBody Copying the body so that the original body with resources can be released -.
func (h HttpRequest) copyBody(src io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(src)

	if err != nil {
		return nil, errors.New("error reading the request body")
	}

	err = src.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// resetBody Resetting in order to close the request body to avoid keeping it open all the time -.
func (h HttpRequest) resetBody(request *http.Request, originalBody []byte) {
	request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
}
