package requests_handler

import (
	attr "Zbolt/pages/request/requests-handler/attributes"
	"net/http"
	"strings"
)

func http_headers_to_attr(http_header http.Header) []attr.Attribute {
	attrs := make([]attr.Attribute, 0, len(http_header))
	for k, v := range http_header {
		val := strings.Join(v, ",")
		if strings.TrimSpace(k) == "" && strings.TrimSpace(val) == "" {
			continue
		}
		attrs = append(attrs, attr.Attribute{
			Key:   k,
			Value: val,
		})
	}
	return attrs
}
