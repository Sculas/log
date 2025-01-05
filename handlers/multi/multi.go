// Package multi implements a handler which invokes a number of handlers.
package multi

import (
	"github.com/sculas/log"
)

// Handler implementation.
type Handler struct {
	Handlers []log.Handler
}

// New handler.
func New(h ...log.Handler) *Handler {
	return &Handler{
		Handlers: h,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	for _, handler := range h.Handlers {
		if err := handler.HandleLog(e); err != nil {
			return err //nolint:wrapcheck // (upstream should handle this)
		}
	}

	return nil
}
