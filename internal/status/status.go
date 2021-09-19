// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package status

import "net/http"

// ClientClosedRequest non-standard HTTP status code for client disconnection.
const ClientClosedRequest = 499

// ClientClosedRequestText non-standard HTTP status for client disconnection.
const ClientClosedRequestText = "Client Closed Request"

func Text(statusCode int) string {
	if statusCode == ClientClosedRequest {
		return ClientClosedRequestText
	}
	return http.StatusText(statusCode)
}
