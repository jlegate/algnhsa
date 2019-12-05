package algnhsa

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
)

const allContentType = "*"

type lambdaResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

func newLambdaResponse(w *httptest.ResponseRecorder, binaryContentTypes map[string]map[string]bool) (lambdaResponse, error) {
	event := lambdaResponse{}

	// Set status code.
	event.StatusCode = w.Code

	// Set headers.
	event.MultiValueHeaders = w.Result().Header

	// Set body.
	fullContentType := w.Header().Get("Content-Type")
	ctParts := strings.Split(fullContentType, "/")
	contentType := ctParts[0]
	contentSubType := ctParts[1]

	fmt.Fprintf(os.Stderr, "Content-Type: %s (%s, %s)\n", fullContentType, contentType, contentSubType)
	if binaryContentTypes[allContentType][allContentType] || binaryContentTypes[contentType][allContentType] || binaryContentTypes[contentType][contentSubType] {
		fmt.Fprintf(os.Stderr, "binary file\n")
		event.Body = base64.RawStdEncoding.EncodeToString(w.Body.Bytes())
		event.IsBase64Encoded = true
	} else {
		fmt.Fprintf(os.Stderr, "non-binary file\n")
		fmt.Fprintf(os.Stderr, "%+v\n", binaryContentTypes)
		event.Body = w.Body.String()
	}

	return event, nil
}
