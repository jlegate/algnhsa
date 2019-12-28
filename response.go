package algnhsa

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"unicode/utf8"
)

const allContentType = "*"

type lambdaResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
	Body              string              `json:"body"`
}

func newLambdaResponse(w *httptest.ResponseRecorder, binaryContentTypes map[string]map[string]bool) (lambdaResponse, error) {
	event := lambdaResponse{}

	// Set status code.
	event.StatusCode = w.Code

	// Set headers.
	event.MultiValueHeaders = w.Result().Header
	bb := w.Body.Bytes()

	// Set body.
	fullContentType := w.Header().Get("Content-Type")
	if fullContentType == "" {
		fullContentType = http.DetectContentType(bb)
	}
	fmt.Fprintf(os.Stderr, "fullContentType: %s\n", fullContentType)
	ctParts := strings.Split(fullContentType, "/")
	fmt.Fprintf(os.Stderr, "%#v\n", ctParts)
	contentType := ctParts[0]
	contentSubType := ctParts[1]

	fmt.Fprintf(os.Stderr, "Content-Type: %s (%s, %s)\n", fullContentType, contentType, contentSubType)

	var output string

	forceBinary := false

	if binaryContentTypes[allContentType][allContentType] || binaryContentTypes[contentType][allContentType] || binaryContentTypes[contentType][contentSubType] {
		forceBinary = true
	}

	if utf8.Valid(bb) || !forceBinary {
		output = string(bb)
		event.IsBase64Encoded = false
	} else {
		output = base64.StdEncoding.EncodeToString(bb)
		event.IsBase64Encoded = true
	}

	event.Body = output

	return event, nil
}
