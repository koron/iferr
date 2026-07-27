package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func iferrStr(in string, pos int, errMsg string) (string, error) {
	out := &bytes.Buffer{}
	r := strings.NewReader(in)
	err := iferr(out, r, pos, errMsg)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func iferrOK(t *testing.T, fn string, off int, errMsg, want string) {
	t.Helper()

	const (
		fnPre   = "package main\nfunc foo() "
		fnPost  = " {}"
		actPre  = "if err != nil {\n\treturn "
		actPost = "\n}\n"
	)

	got, err := iferrStr(fnPre+fn, len(fnPre)+1+off, errMsg)
	if err != nil {
		t.Errorf("iferr() is failed: %s for %q", err, fn)
		return
	}
	if !strings.HasPrefix(got, actPre) || !strings.HasSuffix(got, actPost) {
		t.Errorf("iferr() returns with unexpected prefix or suffix: %q", got)
		return
	}
	got = got[len(actPre) : len(got)-len(actPost)]
	if got != want {
		t.Errorf("iferr() returns unexpected: want=%q got=%q", want, got)
		return
	}
}

func TestIferr(t *testing.T) {
	iferrOK(t, `(interface{}, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(any, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(map[string]struct{}, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(chan bool, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(time.Duration, error)`, 0, `err`, `0, err`)
	iferrOK(t, `(time.Time, error)`, 0, `err`, `time.Time{}, err`)
	iferrOK(t, `(bool, error)`, 0, `err`, `false, err`)
	iferrOK(t, `(foo, error)`, 0, `err`, `foo{}, err`)
	iferrOK(t, `(*foo, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(*foo, error)`, 0, `fmt.Errorf("failed to %v", err)`, `nil, fmt.Errorf("failed to %v", err)`)
}

func TestNumericTypes(t *testing.T) {
	// See https://go.dev/ref/spec#Numeric_types
	for _, typ := range []string{
		"uint", "uint8", "uint16", "uint32", "uint64",
		"int", "int8", "int16", "int32", "int64",
		"float32", "float64",
		"complex64", "complex128",
		"byte",
		"rune",
		"time.Duration",
	} {
		t.Run(typ, func(t *testing.T) {
			iferrOK(t, fmt.Sprintf(`(%s, error)`, typ), 0, `err`, `0, err`)
		})
	}
}
