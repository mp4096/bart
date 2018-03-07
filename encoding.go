package bart

import (
	"encoding/base64"
)

/// Encode a string according to RFC 1342
/// See https://tools.ietf.org/html/rfc1342
func EncodeRfc1342(someString string) string {
	stringInBase64 := base64.StdEncoding.EncodeToString([]byte(someString))
	return "=?utf-8?B?" + stringInBase64 + "?="
}
