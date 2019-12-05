package algnhsa

import (
	"strings"
)

type RequestType int

const (
	RequestTypeAuto RequestType = iota
	RequestTypeAPIGateway
	RequestTypeALB
)

// Options holds the optional parameters.
type Options struct {
	// RequestType sets the expected request type.
	// By default algnhsa deduces the request type from the lambda function payload.
	RequestType RequestType

	// BinaryContentTypes sets content types that should be treated as binary types.
	// The "*/* value makes algnhsa treat any content type as binary.
	BinaryContentTypes   []string
	binaryContentTypeMap map[string]map[string]bool

	// Use API Gateway PathParameters["proxy"] when constructing the request url.
	// Strips the base path mapping when using a custom domain with API Gateway.
	UseProxyPath bool
}

func (opts *Options) setBinaryContentTypeMap() {
	var types = make(map[string]map[string]bool)
	for _, fullContentType := range opts.BinaryContentTypes {
		ctParts := strings.Split(fullContentType, "/")
		var a, b string
		if len(ctParts) != 2 {
			a = fullContentType
			b = "*"
		} else {
			a = ctParts[0]
			b = ctParts[1]
		}
		if _, exists := types[a]; exists != true {
			types[a] = {}
	}
		types[ctParts[0]][ctParts[1]] = true
	}
	opts.binaryContentTypeMap = types
}
