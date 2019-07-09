package idltypes

const(
	ErrorCategoryParse = 1

	ErrorNoDumplicated = 1

	ErrorNoFileTypeUnmatch = 2
	ErrorNoFileTypeInvalid = 3
)


type Error interface {
	error
	Category() int
	Errno() int
}

// kiteErr is the official implementation of KiteError.
type idlParseErr struct {
	category   int
	errno      int
	underlying error
	errMsg     string
}

// Category implements the method from 'KiteError' interface.
func (ke *idlParseErr) Category() int {
	return ke.category
}

// Errno implements the method from 'KiteError' interface.
func (ke *idlParseErr) Errno() int {
	return ke.errno
}

func (ke *idlParseErr)Error() string{
	return ke.errMsg
}

// NewKiteError creates an KiteError instance.
func NewParseError(category int, errno int, underlying error) Error {
	ke := &idlParseErr{category, errno, underlying, ""}
	return ke
}
