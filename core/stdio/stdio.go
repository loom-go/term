package stdio

import (
	"os"
)

var stdin *Reader

func Stdin() *Reader {
	if stdin == nil {
		stdin = NewReader(os.Stdin)
	}

	return stdin
}
