package httprequest

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

// PerformPost used to perform raw post request -.
func (h HttpRequest) PerformPost(args RequestDataParams, backOffStrategy []time.Duration, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, args.Endpoint, bytes.NewBuffer(args.Data))

	if err != nil {
		return nil, err
	}

	// Setting the necessary headers -.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	data, err := h.RetryDo(req, 3, time.Second*10, backOffStrategy)

	if err != nil {
		log.Printf("ERROR: %s EXECUTING POST TO  %s \n", err.Error(), req.URL)
		return nil, err
	}

	return data, nil
}
