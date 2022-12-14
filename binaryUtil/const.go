package binaryUtil

import "errors"

const (
	zero  = byte('0')
	one   = byte('1')
	lsb   = byte('[') // left square brackets
	rsb   = byte(']') // right square brackets
	space = byte(' ')
)

var uint8arr [8]uint8

// ErrBadStringFormat represents a error of input string's format is illegal .
var ErrBadStringFormat = errors.New("bad string format")

// ErrEmptyString represents a error of empty input string.
var ErrEmptyString = errors.New("empty string")

func ErrTypeUnsupport(typeName string) error {
	return errors.New("data type is unsupported: " + typeName)
}
func init() {
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}
