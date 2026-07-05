package jsvm_test

import (
	"testing"

	"github.com/xmx/jsvm"
)

const highJS = `
import foo from 'foo'
foo.bar()
`

func TestTransform(t *testing.T) {
	ret := jsvm.Transform("test", highJS)
	t.Log(string(ret.Code))
}
