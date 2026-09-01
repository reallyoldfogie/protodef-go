package datatypes

import (
	"fmt"
	"os"
)

// DebugEnabled reports whether debug logging is turned on via the
// PROTODEF_GO_DEBUG=1 environment variable.
var DebugEnabled = os.Getenv("PROTODEF_GO_DEBUG") == "1"

// DebugPrintf writes formatted debug output to stderr when PROTODEF_GO_DEBUG=1
// is set. It is a no-op otherwise.
func DebugPrintf(format string, args ...any) {
	if DebugEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

// DebugPrintln writes debug output to stderr when PROTODEF_GO_DEBUG=1 is set.
// It is a no-op otherwise.
func DebugPrintln(args ...any) {
	if DebugEnabled {
		fmt.Fprintln(os.Stderr, args...)
	}
}
