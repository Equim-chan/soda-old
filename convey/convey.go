// package convey provides interfaces for input/output text data in a
// user interactive way.
package convey // import "ekyu.moe/soda/convey"

// TODO: doc: should never copy!
type (
	ReadFunc  func() ([]byte, error)
	WriteFunc func([]byte) error
)

// var (
// 	EditorWrite = editorWrite
// 	EditorRead  = editorWrite

// 	ClipboardWrite = clipboardWrite
// 	ClipboardRead  = clipboardRead

// 	TerminalWrite = terminalWrite
// )

// 为了保持记事本的兼容性，这个函数可以将 LF 换成 CRLF
// func forceCrlf(s string) string {
// 	lfText := strings.Replace(s, "\r\n", "\n", -1)
// 	return strings.Replace(lfText, "\n", "\r\n", -1)
// }

// 为了保持终端的兼容性，这个函数可以将 CRLF 换成 LF
// func forceLf(s string) string {
// 	crlfText := strings.Replace(s, "\n", "\r\n", -1)
// 	return strings.Replace(crlfText, "\r\n", "\n", -1)
// }
