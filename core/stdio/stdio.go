package stdio

import (
	"os"
)

var Stdin = NewReader(os.Stdin)
