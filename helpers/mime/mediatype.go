package mime

import (
	"errors"
	"strings"
)

// This is taken from the built in golang "mime" package
// We wanted a faster ParseMediaType function (we don't care about properties - yet)

// isTSpecial reports whether rune is in 'tspecials' as defined by RFC
// 1521 and RFC 2045.
func isTSpecial(r rune) bool {
	return strings.ContainsRune(`()<>@,;:\"/[]?=`, r)
}

// isTokenChar reports whether rune is in 'token' as defined by RFC
// 1521 and RFC 2045.
func isTokenChar(r rune) bool {
	// token := 1*<any (US-ASCII) CHAR except SPACE, CTLs,
	//             or tspecials>
	return r > 0x20 && r < 0x7f && !isTSpecial(r)
}

func isNotTokenChar(r rune) bool {
	return !isTokenChar(r)
}

// consumeToken consumes a token from the beginning of provided
// string, per RFC 2045 section 5.1 (referenced from 2183), and return
// the token consumed and the rest of the string. Returns ("", v) on
// failure to consume at least one character.
func consumeToken(v string) (token, rest string) {
	notPos := strings.IndexFunc(v, isNotTokenChar)
	if notPos == -1 {
		return v, ""
	}
	if notPos == 0 {
		return "", v
	}
	return v[0:notPos], v[notPos:]
}

func checkMediaTypeDisposition(s string) error {
	typ, rest := consumeToken(s)
	if typ == "" {
		return errors.New("mime: no media type")
	}
	if rest == "" {
		return nil
	}
	if !strings.HasPrefix(rest, "/") {
		return errors.New("mime: expected slash after first token")
	}
	subtype, rest := consumeToken(rest[1:])
	if subtype == "" {
		return errors.New("mime: expected token after slash")
	}
	if rest != "" {
		return errors.New("mime: unexpected content after media subtype")
	}
	return nil
}

// ParseMediaType parses a media type value and any optional
// parameters, per RFC 1521.  Media types are the values in
// Content-Type and Content-Disposition headers (RFC 2183).
// On success, ParseMediaType returns the media type converted
// to lowercase and trimmed of white space
func ParseMediaType(v string) (mediatype string, err error) {
	i := strings.Index(v, ";")
	if i == -1 {
		i = len(v)
	}
	mediatype = strings.TrimSpace(strings.ToLower(v[0:i]))

	err = checkMediaTypeDisposition(mediatype)
	if err != nil {
		return "", err
	}
	return
}
