package main

import (
	"bytes"
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

func iferrOK(t *testing.T, fn string, off int, errMsg, exp string) {
	const (
		fnPre   = "package main\nfunc foo() "
		fnPost  = " {}"
		actPre  = "if err != nil {\n\treturn "
		actPost = "\n}\n"
	)

	act, err := iferrStr(fnPre+fn, len(fnPre)+1+off, errMsg)
	if err != nil {
		t.Errorf("iferr() is failed: %s for %q", err, fn)
		return
	}
	if !strings.HasPrefix(act, actPre) || !strings.HasSuffix(act, actPost) {
		t.Errorf("iferr() returns with unexpected prefix or suffix: %q", act)
		return
	}
	act = act[len(actPre) : len(act)-len(actPost)]
	if act != exp {
		t.Errorf("iferr() returns unexpected: actual=%q expect=%q", act, exp)
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
	iferrOK(t, `(*foo, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(*foo, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(*foo, error)`, 0, `err`, `nil, err`)
	iferrOK(t, `(*foo, error)`, 0, `fmt.Errorf("failed to %v", err)`, `nil, fmt.Errorf("failed to %v", err)`)
}
