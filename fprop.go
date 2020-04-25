package fprop

import (
	"errors"
)

var SkipNode = errors.New("skip this node")

type EncodeFunc func(interface{}) ([]byte, error)

type DecodeFunc func([]byte, interface{}) error

type MapFunc func(string, Tree, interface{}, byte) (Tree, interface{}, byte)

type VisitFunc func(string, interface{}, byte) error
