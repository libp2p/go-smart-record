// Package ir defines the Intermediate Representation (informally, in-memory representation) of smart records.
package ir

import (
	"io"
)

type Node interface {
	WritePretty(w io.Writer) error
}
