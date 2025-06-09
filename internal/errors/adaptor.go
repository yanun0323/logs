package errors

// Error is the interface for supporting errors (github.com/yanun0323/errors)
type Error interface {
	Message() string
	Cause() error
	Stack() []any
	Attributes() []any
}

// Frame is the interface for supporting errors (github.com/yanun0323/errors)
type Frame interface {
	Parameters() (file, function, line string)
}

// Attr is the interface for supporting errors (github.com/yanun0323/errors)
type Attr interface {
	Parameters() (key string, value any)
}
