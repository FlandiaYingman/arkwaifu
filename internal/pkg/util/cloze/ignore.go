package cloze

import "io"

// IgnoreErr closes a closer but ignore its error explicitly.
func IgnoreErr(closer io.Closer) {
	_ = closer.Close()
}
