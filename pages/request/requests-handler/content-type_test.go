package requests_handler_test

import (
	"testing"
	"API-Client/pages/request/requests-handler"
)

var content_type_testcases = []struct {
	content_type requests_handler.ContentType
	t, sub_t string
}{
	{
		content_type: "",
	},
	{
		content_type: "/",
	},
	{
		content_type: "application",
		t: "application",
	},
	{
		content_type: "/;",
	},
	{
		content_type: "application/json",
		t: "application",
		sub_t: "json",
	},
	{
		content_type: "image/png",
		t: "image",
		sub_t: "png",
	},
	{
		content_type: "text/html;charset=UTF-8",
		t: "text",
		sub_t: "html",
	},
}

func TestContentTypeParse(t *testing.T) {
	for i, testcase := range content_type_testcases {
		ct, sub_t := testcase.content_type.Parse()
		if ct != testcase.t {
			t.Fatalf("testcase %v: Expected %s but got %s\n", i, testcase.t, ct)
		}
		
		if sub_t != testcase.sub_t {
			t.Fatalf("testcase %v: Expected %s but got %s\n", i, testcase.sub_t, ct)
		}
	}
}
