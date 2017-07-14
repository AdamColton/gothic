package lensnippet

import (
	"fmt"
	"github.com/adamcolton/gothic"
)

const (
	minMode byte = iota
	maxMode
	equalMode
)

type lengthSnippet struct {
	length     int
	errSnippet gothic.Snippet
	on         string
	mode       byte
}

func Min(length int, errSnippet gothic.Snippet) gothic.Snippet {
	return lengthSnippet{
		length:     length,
		errSnippet: errSnippet,
		mode:       minMode,
	}
}

func Max(length int, errSnippet gothic.Snippet) gothic.Snippet {
	return lengthSnippet{
		length:     length,
		errSnippet: errSnippet,
		mode:       maxMode,
	}
}

func (ls lengthSnippet) AddContext(key, value string) gothic.Snippet {
	if key == "On" {
		ls.on = value
	}
	return ls
}

const (
	minErrStr = `%s must be at least %d long`
	minStr    = `if len(%s) < %d {
	%s
}`
	maxErrStr = `%s must be no more than %d long`
	maxStr    = `if len(%s) > %d {
	%s
}`
	eqErrStr = `%s must be %d long`
	eqStr    = `if len(%s) != %d {
	%s
}`
)

func (ls lengthSnippet) String() string {
	switch ls.mode {
	case minMode:
		es := ls.errSnippet.AddContext("Error", fmt.Sprintf(minErrStr, ls.on, ls.length))
		return fmt.Sprintf(minStr, ls.on, ls.length, es)
	case maxMode:
		es := ls.errSnippet.AddContext("Error", fmt.Sprintf(maxErrStr, ls.on, ls.length))
		return fmt.Sprintf(maxStr, ls.on, ls.length, es)
	case equalMode:
		es := ls.errSnippet.AddContext("Error", fmt.Sprintf(eqErrStr, ls.on, ls.length))
		return fmt.Sprintf(eqStr, ls.on, ls.length, es)
	}
	return ""
}
