// package convey provides interfaces for input/output text data in a
// user interactive way.
package convey // import "ekyu.moe/soda/convey"

// TODO: doc: should never copy!
type (
	ReadFunc  func() ([]byte, error)
	WriteFunc func([]byte) error
)
