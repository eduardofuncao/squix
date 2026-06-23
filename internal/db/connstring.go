package db

import (
	"net/url"
	"strings"
)

// EncodeUserinfo percent-encodes reserved characters in the username and
// password of URL-format connection strings. Go's net/url.Parse treats reserved
// characters like '#' as structural delimiters, which truncates userinfo and
// breaks DSNs whose passwords contain them.
//
// The password is treated as a literal: any existing '%XX' sequences are
// re-encoded (so '%23' in the config reaches the DB as the literal "%23", not
// decoded to '#').
//
// Strings without "://" (keyword-format DSNs, file paths) are returned
// unchanged.
func EncodeUserinfo(s string) string {
	schemeIdx := strings.Index(s, "://")
	if schemeIdx < 0 {
		return s
	}
	userinfoStart := schemeIdx + 3

	// The userinfo delimiter is the first '@' after the scheme. Using the last
	// '@' would mistake an '@' in the path/query/fragment (e.g. ?role=a@b) for
	// the delimiter and re-encode across the host, corrupting the DSN.
	relAt := strings.Index(s[userinfoStart:], "@")
	if relAt < 0 {
		return s
	}
	atIdx := userinfoStart + relAt

	user, password, found := strings.Cut(s[userinfoStart:atIdx], ":")
	if !found {
		return s
	}

	return s[:userinfoStart] + url.UserPassword(user, password).String() + s[atIdx:]
}
