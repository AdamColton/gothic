package gothicgo

import (
	"fmt"
)

type errStr string

func (err errStr) Error() string {
	return string(err)
}

// ErrCtx helps provide context to errors
type ErrCtx struct {
	Base   error
	ErrStr string
}

func (err ErrCtx) Error() string {
	return fmt.Sprintf("%s\n%v", err.ErrStr, err.Base)
}

func errCtx(err error, format string, data ...interface{}) error {
	if err == nil {
		return nil
	}
	ctxStr := fmt.Sprintf(format, data...)
	if ec, ok := err.(ErrCtx); ok {
		return ErrCtx{
			Base:   ec.Base,
			ErrStr: fmt.Sprintf("%s\n%s", ctxStr, ec.ErrStr),
		}
	}
	return ErrCtx{
		Base:   err,
		ErrStr: ctxStr,
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
