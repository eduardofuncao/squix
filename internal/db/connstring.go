package db

import (
	"net/url"
	"strings"
)

// EncodeUserinfo percent-encodes reserved characters in the username and
// password of URL-format connection strings. Go's net/url.Parse treats reserved
// characters like '#' as structural delimiters, which truncates userinfo and
// breaks DSNs whose credentials contain them.
//
// Credentials are treated as literals: any existing '%XX' sequences are
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

	atIdx, ok := findUserinfoDelimiter(s, userinfoStart)
	if !ok {
		return s
	}

	user, password, found := strings.Cut(s[userinfoStart:atIdx], ":")
	if !found {
		// Username-only userinfo: encode the username alone.
		return s[:userinfoStart] + url.User(user).String() + s[atIdx:]
	}

	return s[:userinfoStart] + url.UserPassword(user, password).String() + s[atIdx:]
}

// findUserinfoDelimiter returns the index of the '@' that separates userinfo
// from host. It is the first '@' whose following segment (up to the next '/',
// '?' or '#') contains no other '@' — i.e. a plausible host. This lets a
// password containing '@' be encoded (%40) instead of truncating the userinfo
// early, while an '@' in the path/query/fragment is never mistaken for the
// delimiter.
func findUserinfoDelimiter(s string, userinfoStart int) (int, bool) {
	from := userinfoStart
	for {
		rel := strings.Index(s[from:], "@")
		if rel < 0 {
			return 0, false
		}
		atIdx := from + rel

		hostSeg := s[atIdx+1:]
		if end := strings.IndexAny(hostSeg, "/?#"); end >= 0 {
			hostSeg = hostSeg[:end]
		}
		if !strings.Contains(hostSeg, "@") {
			return atIdx, true
		}
		from = atIdx + 1
	}
}
