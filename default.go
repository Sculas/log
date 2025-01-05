package log

import (
	"bytes"
	"fmt"
	"log"
)

// handleStdLog outputs to the stdlib log.
func handleStdLog(e *Entry) error {
	level := levelNames[e.Level]
	fields := e.Fields.Sorted()

	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, "%5s %-25s", level, e.Message)
	for _, pair := range fields {
		_, _ = fmt.Fprintf(&b, " %s=%v", pair.Key, pair.Value)
	}
	log.Println(b.String())

	return nil
}
