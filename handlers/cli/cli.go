// Package cli implements a tracing-like logging format.
package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/mattn/go-colorable"
	"github.com/rotisserie/eris"
	"github.com/sculas/log"
)

var (
	Default    = New(os.Stderr, true)
	colorField = color.New(color.OpItalic)
	colorMuted = color.New(color.FgDarkGray)
	colorError = color.New(color.Bold, color.FgRed)
)

// Colors mapping.
var Colors = [...]color.Style{
	log.DebugLevel: color.New(color.Bold, color.FgBlue),
	log.InfoLevel:  color.New(color.Bold, color.FgGreen),
	log.WarnLevel:  color.New(color.Bold, color.FgYellow),
	log.ErrorLevel: color.New(color.Bold, color.FgRed),
	log.FatalLevel: color.New(color.Bold, color.FgRed),
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  " INFO",
	log.WarnLevel:  " WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

type Handler struct {
	mu      sync.Mutex
	errfmt  eris.StringFormat
	Writer  io.Writer
	Padding int
}

func New(w io.Writer, useColors bool) *Handler {
	errfmt := eris.NewDefaultStringFormat(eris.FormatOptions{
		WithTrace:    true, // enable stack trace
		InvertTrace:  true, // invert stack trace (top to bottom)
		WithExternal: true, // format external errors not from eris
	})

	if f, ok := w.(*os.File); ok {
		if useColors {
			return &Handler{errfmt: errfmt, Writer: colorable.NewColorable(f), Padding: 2}
		}
	}

	return &Handler{errfmt: errfmt, Writer: colorable.NewNonColorable(w), Padding: 2}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	level, levelColor := Strings[e.Level], Colors[e.Level]
	fields := e.Fields.Sorted()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Write the timestamp and level.
	_, _ = fmt.Fprint(h.Writer,
		colorMuted.Render(time.Now().Format("2006-01-02T15:04:05.000000Z")),
		" ",
		levelColor.Render(level),
		" ")

	// Write the logger name, if present.
	if e.Name != "" {
		_, _ = fmt.Fprint(h.Writer, colorMuted.Render(e.Name+": "))
	}

	// Write the log message, if present.
	if e.Message != "" {
		_, _ = fmt.Fprint(h.Writer, e.Message+" ")
	}

	// Write the fields (except the name of the logger).
	for _, field := range fields {
		if field.Value == nil {
			field.Value = ""
		}
		_, _ = fmt.Fprint(h.Writer,
			colorField.Render(field.Key),
			colorMuted.Render("="),
			fmt.Sprintf(`"%v"`, field.Value),
			" ")
	}

	// Write the new line (as we haven't done that yet).
	_, _ = fmt.Fprintln(h.Writer)

	// Write the error if there is one, but only in debug mode.
	if e.Logger.Level == log.DebugLevel {
		for _, field := range fields {
			if field.Key == "error" {
				if err, ok := field.Value.(error); ok {
					_, _ = fmt.Fprintf(h.Writer, "\n%s\n%+v\n\n",
						colorError.Render("Stacktrace:"),
						eris.ToCustomString(err, h.errfmt))
				}

				break
			}
		}
	}

	return nil
}
