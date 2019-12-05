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

	// Set body.
	fullContentType := w.Header().Get("Content-Type")
	ctParts := strings.Split(fullContentType, "/")
	contentType := ctParts[0]
	contentSubType := ctParts[1]

	fmt.Fprintf(os.Stderr, "Content-Type: %s (%s, %s)\n", fullContentType, contentType, contentSubType)

	var output string

	bb := w.Body.Bytes()

	fmt.Fprintf(os.Stderr, "Bytes: %+v\n", string(bb))

	if utf8.Valid(bb) {
		fmt.Fprintf(os.Stderr, "non-binary file\n")
		output = string(bb)
		event.IsBase64Encoded = false
	} else {
		fmt.Fprintf(os.Stderr, "binary file\n")
		output = base64.StdEncoding.EncodeToString(bb)
		event.IsBase64Encoded = true
	}

	// if binaryContentTypes[allContentType][allContentType] || binaryContentTypes[contentType][allContentType] || binaryContentTypes[contentType][contentSubType] {
	// 	fmt.Fprintf(os.Stderr, "binary file\n")
	// 	event.Body = base64.RawStdEncoding.EncodeToString(w.Body.Bytes())
	// 	event.IsBase64Encoded = true
	// } else {
	// 	fmt.Fprintf(os.Stderr, "non-binary file\n")
	// 	fmt.Fprintf(os.Stderr, "%+v\n", binaryContentTypes)
	// 	event.Body = w.Body.String()
	// }
	fmt.Fprintf(os.Stderr, "%+v\n", event)
	event.Body = output

	return event, nil
}
