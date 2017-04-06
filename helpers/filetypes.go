package helpers

import "mime"

const (
	// JSMimeType is the mime type for javascript files.
	JSMimeType = "text/javascript"
	// CSSMimeType is the mime type for CSS files.
	CSSMimeType = "text/css"
	// TXTMimeType is the mime type for Text files
	TXTMimeType = "text/plain"
	// JSONMimeType is the mime type for JSon files.
	JSONMimeType = "text/json"
)

func init() {
	mime.AddExtensionType(".js", JSMimeType)
	mime.AddExtensionType(".css", CSSMimeType)
	mime.AddExtensionType(".txt", TXTMimeType)
	mime.AddExtensionType(".json", JSONMimeType)
}
