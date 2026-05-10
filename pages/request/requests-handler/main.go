package requests_handler

import (
	attr "API-Client/pages/request/requests-handler/attributes"
	"path/filepath"

	gui "github.com/guigui-gui/guigui"
)

type RequestType uint8

const (
	HTTP RequestType = iota + 0
	Websocket
	GraphQL
	Grpc
)

func (t RequestType) IconName() string {
	switch t {
	case HTTP:
		return "large-icons/http"
	case Websocket:
		return "large-icons/websocket"
	case GraphQL:
		return "large-icons/graphql"
	case Grpc:
		return "large-icons/grpc"
	default:
		panic("Unknown request type")
	}
}

type Request struct {
	Type RequestType
	path string
	data any // pointer to data
}

func (r *Request) Data() any {
	return r.data
}

func (r *Request) Path() string {
	return r.path
}

func (r *Request) Name() string {
	return filepath.Base(r.path)
}

// Clear deletes the data in RAM
func (r *Request) Clear() {

}

func NewRequest(t RequestType, path string) Request {
	req := Request{
		Type: t,
		path: path,
	}
	if t == HTTP {
		data := HTTP_Data{}
		data.ResponseConfig.AutoWrap = true
		data.ResponseConfig.Formate = true
		data.Headers = []attr.AttrCheck{
			{
				Checked: true,
				Key:     "Accept",
				Value:   "*/*",
			},
			{
				Checked: false,
				Key:     "Accept-Encoding",
				Value:   "gzip, deflate, br, zstd",
			},
			{
				Checked: true,
				Key:     "User-Agent",
				Value:   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
			},
		}
		req.data = &data
	}
	return req
}

type RequestWidget interface {
	gui.Widget
	RequestType() RequestType
	SetReq(req *Request)
	SyncData()
}

type Folder struct {
	path string
}

func (r *Folder) Path() string {
	return r.path
}

func (r *Folder) Name() string {
	return filepath.Base(r.path)
}

func NewFolder(path, name string) Folder {
	return Folder{
		path: filepath.Join(path, name),
	}
}

type FolderOrFile interface {
	Path() string
	Name() string
}
