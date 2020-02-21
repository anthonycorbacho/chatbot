package bot

// Bot errors.
const (
	ErrEmptyPattern = Error("empty regexp expression")
)

// Error represents a session error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}
