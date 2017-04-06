package helpers

import "mime"

const (
	JSMimeType   = "text/javascript"
	CSSMimeType  = "text/css"
	TXTMimeType  = "text/plain"
	JSONMimeType = "text/json"
)

func init() {
	mime.AddExtensionType(".js", JSMimeType)
	mime.AddExtensionType(".css", CSSMimeType)
	mime.AddExtensionType(".txt", TXTMimeType)
	mime.AddExtensionType(".json", JSONMimeType)
}
