package url_utils

import "net/url"

// CleanURL removes raw query and fragments.
func CleanURL(u *url.URL) {
	if u == nil {
		return
	}
	u.RawQuery = ""
	u.RawFragment = ""
	u.Fragment = ""
	u.ForceQuery = false
}
