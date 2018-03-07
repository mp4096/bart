package bart

import (
	"bytes"
	"encoding/base64"
)

// EncodeRfc1342 encodes a string according to RFC 1342
// See https://tools.ietf.org/html/rfc1342
func EncodeRfc1342(someString string) string {
	stringInBase64 := base64.StdEncoding.EncodeToString([]byte(someString))
	return "=?utf-8?B?" + stringInBase64 + "?="
}

// EscapeAngleBrackets escapes angle brackets for HTML
func EscapeAngleBrackets(someBytes []byte) []byte {
	escaped := bytes.Replace(someBytes, []byte("<"), []byte("&lt;"), -1)
	return bytes.Replace(escaped, []byte(">"), []byte("&gt;"), -1)
}
