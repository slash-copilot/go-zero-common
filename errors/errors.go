package errors

import "fmt"

// CodeMsg is a struct that contains a code and a message.
// It implements the error interface.
type CodeMsg struct {
	Code string
	Msg  string
}

func (c *CodeMsg) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", c.Code, c.Msg)
}

// New creates a new CodeMsg.
func New(code string, msg string) error {
	return &CodeMsg{Code: code, Msg: msg}
}
