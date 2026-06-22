package requests_handler

import (
	"mime"
	"slices"
	"strings"
)

type ContentType string

// These are common content types with thir most populer extension name.
// This is here becuase text/html like content types returns mhtml like extension names which are not populer
// and users wouldn't now what is the meaning with those extension name.
var common_content_types = map[string]string{
	"text/html":        "html",
	"text/css":         "css",
	"text/javascript":  "js",
	"application/json": "json",
	"image/png":        "png",
	"image/jpeg":       "jpg",
	"image/gif":        "gif",
	"image/svg+xml":    "svg",
	"image/webp":       "webp",
	"application/pdf":  "pdf",
	"text/csv":         "csv",
	"text/plain":       "txt",
	"application/xml":  "xml",
	"video/mp4":        "mp4",
	"audio/mpeg":       "mp3",
	"application/zip":  "zip",
}

func (content_type ContentType) Extension() string {
	c_type, _, _ := strings.Cut(string(content_type), ";")
	c_type, ok := common_content_types[c_type]
	if ok {
		return c_type
	}

	ex, _ := mime.ExtensionsByType(string(content_type))
	if len(ex) == 0 {
		_, sub_t := content_type.Parse()
		return strings.ToLower(sub_t)
	}
	return strings.TrimPrefix(ex[0], ".")
}

func (content_type ContentType) Extensions() []string {
	best_extension := content_type.Extension()
	extensions, _ := mime.ExtensionsByType(string(content_type))
	for i, ex := range extensions {
		extensions[i] = strings.TrimPrefix(ex, ".")
	}

	extensions = append([]string{best_extension}, extensions...)

	i := slices.Index(extensions, best_extension)
	if i != -1 {
		extensions = slices.Delete(extensions, i, i+1)
	}

	return extensions
}

func (content_type ContentType) Parse() (t, sub_t string) {
	if content_type == "" {
		return
	}

	var i int
	for ; i < len(content_type) && content_type[i] != '/'; i++ {
	}

	if i == len(content_type) {
		t = string(content_type)
		return
	} else if content_type[i] == '/' {
		t = string(content_type[:i])
	}

	j := i
	for ; i < len(content_type) && content_type[i] != ';'; i++ {
	}

	if i == len(content_type) {
		sub_t = string(content_type[j+1 : i])
	} else if content_type[i] == ';' {
		sub_t = string(content_type[j+1 : i])
	}

	return
}
