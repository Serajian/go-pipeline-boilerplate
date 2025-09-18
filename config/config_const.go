package config

// ******* PATH *******

const (
	// DefaultCFGPath default path for config file
	DefaultCFGPath string = "./config/config.example.yaml"
)

// *******CONST in Main*******
type contextKey int

const (
	// TraceIDKey key for traceID
	TraceIDKey contextKey = iota
)

// *******Channels Stages*******

const (
	BuffData       int = 64
	BuffErr        int = 64
	BuffBarrierCap int = 16
)
