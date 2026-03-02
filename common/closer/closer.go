// Package closer provides utilities for closing resources with error logging.
package closer

import (
	"errors"
	"io"
	"net"
)

// LogFunc is a logging function compatible with log.Printf.
type LogFunc func(string, ...interface{})

// LogClose closes c and logs any error using the provided logger.
// Errors from closing an already-closed network connection are silently ignored.
func LogClose(c io.Closer, logf LogFunc) {
	if err := c.Close(); err != nil {
		var netErr *net.OpError
		if errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed) {
			return
		}
		logf("close error: %v", err)
	}
}
