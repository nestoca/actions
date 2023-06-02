package logging

import (
	"fmt"
	"os"

	"github.com/TwinProduction/go-color"
)

func Log(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, color.Ize(color.Yellow, format+"\n"), a...)
}
