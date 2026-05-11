package url_utils_test

import (
	url_utils "Zbolt/pages/request/requests-handler/url-utils"
	"testing"
)

func TestParse_url_path_query(t *testing.T){
	pattern, _ := url_utils.ParsePattern("/{user}/{repo}")
	
	_, ok := pattern.Values["user"]
	if !ok {
		t.Fatal("Expected user but got", ok)
	}
	pattern.Set("user", "udan-jayanith")
	
	_, ok = pattern.Values["repo"]
	if !ok {
		t.Fatal("Expected repo but got", ok)
	}
	pattern.Set("repo", "zbolt")
	
	if pattern.Path() != "/udan-jayanith/zbolt" {
		t.Fatal("Unexpected output", pattern.Path())
	}
}

