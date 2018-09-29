package bart

import (
	"bytes"
	"encoding/base64"
)

// EncodeToBase64 transforms a string to bytes and encodes with Base64
func EncodeToBase64(someString string) string {
	return base64.StdEncoding.EncodeToString([]byte(someString))
}

// EncodeToBase64WrapLines encodes a string with Base64 and inserts linebreaks
// after every 76 characters
func EncodeToBase64WrapLines(someString string) []byte {
	const maxLineLength = 76

	bytes := []byte(EncodeToBase64(someString))
	numLines := len(bytes) / maxLineLength
	// 2 bytes for each CRLF at the line end
	wrappedBytes := make([]byte, 0, len(bytes)+2*numLines)
	for i := 0; i < numLines; i++ {
		wrappedBytes = append(wrappedBytes, bytes[i*maxLineLength:(i+1)*maxLineLength]...)
		wrappedBytes = append(wrappedBytes, []byte("\r\n")...)
	}
	wrappedBytes = append(wrappedBytes, bytes[numLines*maxLineLength:]...)
	wrappedBytes = append(wrappedBytes, []byte("\r\n")...)
	return wrappedBytes
}

// EncodeRfc1342 encodes a string according to RFC 1342
// See https://tools.ietf.org/html/rfc1342
func EncodeRfc1342(someString string) string {
	return "=?utf-8?B?" + EncodeToBase64(someString) + "?="
}

// EscapeHtmlCharacters escapes angle brackets and ampersand for HTML
func EscapeHtmlCharacters(someBytes []byte) []byte {
	escaped := bytes.Replace(someBytes, []byte("&"), []byte("&amp;"), -1)
	escaped = bytes.Replace(escaped, []byte("<"), []byte("&lt;"), -1)
	return bytes.Replace(escaped, []byte(">"), []byte("&gt;"), -1)
}
