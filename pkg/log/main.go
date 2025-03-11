package log

// ColorEnabled returns true if color is enabled for the logger
// this should be used to control output
func ColorEnabled(logger Logger) bool {
	type maybeColorer interface {
		ColorEnabled() bool
	}
	v, ok := logger.(maybeColorer)
	return ok && v.ColorEnabled()
}
