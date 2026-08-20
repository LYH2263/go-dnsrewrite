package errors

import "fmt"

// WrapErr 用 %w 包装 sentinel 与底层错误。
func WrapErr(sentinel, err error) error {
	if sentinel == nil {
		return err
	}
	if err == nil {
		return sentinel
	}
	return fmt.Errorf("%w: %w", sentinel, err)
}

// WrapMsg 包装带消息。
func WrapMsg(sentinel error, msg string) error {
	if sentinel == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%w: %s", sentinel, msg)
}
