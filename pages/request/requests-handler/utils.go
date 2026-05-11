package requests_handler

import (
	attr "Zbolt/pages/request/requests-handler/attributes"
	"net/http"
	"strings"
)

func http_headers_to_attr_check(http_header http.Header) []attr.AttrCheck {
	attrs := make([]attr.AttrCheck, 0, len(http_header))
	for k, v := range http_header {
		attrs = append(attrs, attr.AttrCheck{
			Key:     k,
			Value:   strings.Join(v, ","),
			Checked: true,
		})
	}
	return attrs
}
