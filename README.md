# Generate "if err != nil {" block

Generate `if err != nil {` block for current function.

## Usage

Install and update by

```console
$ go get -u github.com/koron/iferr
```

Run, it get `if err != nil {` block for the postion at 1234 bytes.

```console
$ iferr -pos 1234 < main.go
if err != nil {
	return ""
}
```
