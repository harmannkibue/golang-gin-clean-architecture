package httprequest

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

// Helper function to compare byte slices -.
func stringSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Mock ReadCloser implementation for testing -.
type mockReadCloser struct {
	err error
}

func (m mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func (m mockReadCloser) Close() error {
	return m.err
}

func TestRetryDo(t *testing.T) {
	// Create a mock request for testing -.
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	h := HttpRequest{}

	// Test case: successful request -.
	mockClient := &mockHTTPClient{}
	mockClient.resp = &http.Response{
		StatusCode: http.StatusOK,
	}
	_, err = h.RetryDo(req, 3, time.Second, []time.Duration{time.Millisecond, time.Millisecond, time.Millisecond})
	if err != nil {
		t.Errorf("RetryDo returned an unexpected error: %v", err)
	}

	// Test case: request fails with status code 400
	mockClient.resp = &http.Response{
		StatusCode: http.StatusBadRequest,
	}

	// Simulate a 404 request with the wrong path to test domain -.
	req404, err := http.NewRequest("GET", "https://example.com/errrr", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = h.RetryDo(req404, 3, time.Second, []time.Duration{time.Millisecond, time.Millisecond, time.Millisecond})

	// checking if the error returned is not 404 -.
	if err != nil && err.Error() != "failed with status code 404" {
		t.Errorf("RetryDo returned an unexpected error: %v", err)
	}

}

// Testing resetting the body during http request -.
func TestResetBody(t *testing.T) {
	testCases := []struct {
		name          string
		originalBody  []byte
		expectedBody  []byte
		expectedError error
	}{
		{
			name:          "Valid body",
			originalBody:  []byte("Valid request body!"),
			expectedBody:  []byte("Valid request body!"),
			expectedError: nil,
		},
		{
			name:          "Empty body",
			originalBody:  []byte{},
			expectedBody:  []byte{},
			expectedError: nil,
		},
	}

	h := HttpRequest{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &http.Request{}
			h.resetBody(request, tc.originalBody)

			// Check if the body was reset correctly
			body, err := io.ReadAll(request.Body)
			if err != nil {
				t.Errorf("Error reading request body: %s", err)
			}
			if !bytes.Equal(body, tc.expectedBody) {
				t.Errorf("Unexpected body. Got: %s, Expected: %s", body, tc.expectedBody)
			}

			// Check if GetBody function returns the expected body
			getBody, err := request.GetBody()
			if err != nil {
				t.Errorf("Error calling GetBody: %s", err)
			}
			readBody, err := io.ReadAll(getBody)
			if err != nil {
				t.Errorf("Error reading GetBody: %s", err)
			}
			if !bytes.Equal(readBody, tc.expectedBody) {
				t.Errorf("Unexpected GetBody. Got: %s, Expected: %s", readBody, tc.expectedBody)
			}
		})
	}
}

// Testing copying httprequest body -.
func TestCopyBody(t *testing.T) {
	testCases := []struct {
		name          string
		inputBody     io.ReadCloser
		expectedBytes []byte
		expectedError error
	}{
		{
			name:          "Valid body",
			inputBody:     ioutil.NopCloser(strings.NewReader("Valid request body!")),
			expectedBytes: []byte("Valid request body!"),
			expectedError: nil,
		},
		{
			name:          "Empty body",
			inputBody:     ioutil.NopCloser(strings.NewReader("")),
			expectedBytes: []byte{},
			expectedError: nil,
		},
		{
			name:          "Error reading body",
			inputBody:     mockReadCloser{err: errors.New("read error")},
			expectedBytes: nil,
			expectedError: errors.New("error reading the request body"),
		},
		{
			name:          "Error closing body",
			inputBody:     mockReadCloser{err: errors.New("close error")},
			expectedBytes: nil,
			expectedError: errors.New("error reading the request body"),
		},
	}

	h := HttpRequest{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := h.copyBody(tc.inputBody)

			// Check if the returned bytes match the expected result
			if !stringSliceEqual(bytes, tc.expectedBytes) {
				t.Errorf("Unexpected bytes. Got: %s, Expected: %s", bytes, tc.expectedBytes)
			}

			// Check if the returned error matches the expected error
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) || (err != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Unexpected error. Got: %v, Expected: %v", err, tc.expectedError)
			}
		})
	}
}
