package requests_handler

import (
	attr "Zbolt/pages/request/requests-handler/attributes"
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

func (r *Request) Name() string {
	return filepath.Base(r.path)
}

func (r *Request) Path() string {
	return r.path
}

func (r *Request) Rename(name string) error {
	path := filepath.Dir(r.path)
	r.path = filepath.Join(path, name)
	return nil
}

// Clear deletes the data in RAM
func (r *Request) Close() error {
	return nil
}

func NewRequest(t RequestType, path, name string) *Request {
	req := Request{
		Type: t,
		path: filepath.Join(path, name),
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
	return &req
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

func (r *Folder) Rename(name string) error {
	path := filepath.Dir(r.path)
	r.path = filepath.Join(path, name)
	return nil
}

func NewFolder(path, name string) *Folder {
	return &Folder{
		path: filepath.Join(path, name),
	}
}

type Item interface {
	Path() string
	Name() string
	Rename(name string) error
}

type RequestWidget interface {
	gui.Widget
	RequestType() RequestType
	SetReq(req *Request)
	SyncData()
	Close() error
}
