package url_utils

import (
	attr "API-Client/pages/request/requests-handler/attributes"
	"fmt"
	"net/url"
	"strings"
)

// Adapted from Golangs net/url package
func ParseParametersAsCheck(encoded_parameters string) ([]attr.AttrCheck, error) {
	var err error
	list := make([]attr.AttrCheck, 0, 2)
	for encoded_parameters != "" {
		var key string
		key, encoded_parameters, _ = strings.Cut(encoded_parameters, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}

		key, value, _ := strings.Cut(key, "=")
		key, err = url.QueryUnescape(key)
		if err != nil {
			continue
		}

		value, err = url.QueryUnescape(value)
		if err != nil {
			continue
		}

		list = append(list, attr.AttrCheck{
			Checked: true,
			Key:     key,
			Value:   value,
		})
	}

	return list, err
}

func EncodeParameters(params []attr.AttrCheck) string {
	if len(params) == 0 {
		return ""
	}
	var buf strings.Builder
	// This assumes key and values is about length of 5 each.
	buf.Grow(len(params) * 10)

	for _, attr := range params {
		if !attr.Checked {
			continue
		}

		key := url.QueryEscape(attr.Key)
		value := url.QueryEscape(attr.Value)

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(value)
	}
	return buf.String()
}
