package db

import (
	"fmt"
	"strings"
)

// EncodeUserinfoPassword percent-encodes reserved characters in the password
// portion of URL-format connection strings. Go's net/url.Parse treats '#' as
// a fragment delimiter, which truncates userinfo and breaks DSNs whose
// passwords contain it. Already-encoded %XX sequences are preserved.
//
// Strings without "://" (keyword-format DSNs, file paths) are returned
// unchanged.
func EncodeUserinfoPassword(s string) string {
	schemeIdx := strings.Index(s, "://")
	if schemeIdx < 0 {
		return s
	}
	userinfoStart := schemeIdx + 3

	lastAt := strings.LastIndex(s, "@")
	if lastAt < userinfoStart {
		return s
	}

	userinfo := s[userinfoStart:lastAt]
	user, password, found := strings.Cut(userinfo, ":")
	if !found {
		return s
	}

	return s[:userinfoStart] + user + ":" + encodePassword(password) + s[lastAt:]
}

func encodePassword(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		if c == '%' && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
			b.WriteString(s[i : i+3])
			i += 3
			continue
		}
		if shouldEncode(c) {
			fmt.Fprintf(&b, "%%%02X", c)
		} else {
			b.WriteByte(c)
		}
		i++
	}
	return b.String()
}

func shouldEncode(c byte) bool {
	switch c {
	case '#', '?', '/', ' ', '"', '\'':
		return true
	}
	return false
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
