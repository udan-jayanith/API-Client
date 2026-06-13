package url_utils

import "net/url"

// CleanURL removes raw query and fragments.
func CleanURL(u *url.URL) {
	u.RawQuery = ""
	u.RawFragment = ""
	u.Fragment = ""
	u.ForceQuery = false
}

func IsJustPortNumber(u_str string) bool {
	if len(u_str) <= 1 {
		return false
	} else if u_str[0] != ':' {
		return false
	}
	u_str = u_str[1:]
	var i int
	for _, char := range u_str {
		if char < '0' || char > '9' {
			break
		}
		i++
	}
	if i > 0 {
		return true
	}

	return false
}
