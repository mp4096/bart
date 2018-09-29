package bart

const boundary = "=_=_=dKJUIzdDfdjILoJdsOqQYYXSs4RXbH0WgUL"
const delimiter = "\r\n--" + boundary + "\r\n"
const encodingBase64 = "Content-Transfer-Encoding: base64\r\n"
const closingDelimiter = "\r\n--" + boundary + "--\r\n"

type emailHeader struct {
	content string
}

func (h *emailHeader) asPlainBytes() []byte {
	return []byte(h.content)
}

func (h *emailHeader) asHtml() []byte {
	res := []byte("<!DOCTYPE html>\r\n" +
		"<html><head><meta charset=\"UTF-8\"></head>\r\n" +
		"<pre style=\"background-color: #f0f8ff;\">\r\n")
	res = append(res, EscapeHtmlCharacters(h.asPlainBytes())...)
	return append(res, []byte("</pre>\r\n")...)
}

func footerAsPlainBytes() []byte {
	return []byte(closingDelimiter)
}

func footerAsHtml() []byte {
	return []byte("<pre style=\"background-color: #f0f8ff;\">" +
		closingDelimiter +
		"</pre></html>\r\n")
}

type mimepart interface {
	asPlainBytes() []byte
	asHtml() []byte
	asBase64() []byte
	contentType() string
}

type textHtml struct {
	content string
}

func (p *textHtml) asPlainBytes() []byte {
	return []byte(delimiter + p.contentType() + p.content)
}

func (p *textHtml) asBase64() []byte {
	res := []byte(delimiter + encodingBase64 + p.contentType())
	return append(res, EncodeToBase64WrapLines(p.content)...)
}

func (p *textHtml) asHtml() []byte {
	return []byte("<pre style=\"background-color: #fff0f5;\">" +
		delimiter +
		p.contentType() +
		"</pre>\r\n" +
		p.content)
}

func (*textHtml) contentType() string {
	return "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
}

type textPlain struct {
	content string
}

func (p *textPlain) asPlainBytes() []byte {
	return []byte(delimiter + p.contentType() + p.content)
}

func (p *textPlain) asBase64() []byte {
	res := []byte(delimiter + encodingBase64 + p.contentType())
	return append(res, EncodeToBase64WrapLines(p.content)...)
}

func (p *textPlain) asHtml() []byte {
	res := []byte("<pre style=\"background-color: #ffefd5;\">" +
		delimiter +
		p.contentType() +
		"</pre>\r\n<pre>\r\n")
	res = append(res, EscapeHtmlCharacters([]byte(p.content))...)
	return append(res, []byte("</pre>\r\n")...)
}

func (*textPlain) contentType() string {
	return "Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"
}
