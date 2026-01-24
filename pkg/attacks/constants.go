package attacks

import "time"

// Constants for Layer 4 attacks
const (
	// DefaultPacketSize is the default size of data packets sent in TCP/UDP attacks
	DefaultPacketSize = 1024

	// DefaultUDPIterations is the default number of UDP packets sent per connection
	DefaultUDPIterations = 100

	// DefaultSYNIterations is the default number of SYN attempts per execution
	DefaultSYNIterations = 10

	// DefaultDialTimeout is the default timeout for establishing TCP connections
	DefaultDialTimeout = 1 * time.Second

	// SYNDialTimeout is the timeout for SYN flood connection attempts
	SYNDialTimeout = 100 * time.Millisecond
)

// Constants for Layer 7 attacks
const (
	// WorkerTickInterval is the interval at which the work producer generates tasks
	WorkerTickInterval = 1 * time.Millisecond

	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 10 * time.Second

	// DefaultTLSHandshakeTimeout is the default timeout for TLS handshakes
	DefaultTLSHandshakeTimeout = 5 * time.Second

	// DefaultKeepAliveInterval is the default keep-alive interval for connections
	DefaultKeepAliveInterval = 30 * time.Second

	// DefaultMaxIdleConns is the default maximum number of idle connections
	DefaultMaxIdleConns = 100

	// DefaultMaxIdleConnsPerHost is the default maximum idle connections per host
	DefaultMaxIdleConnsPerHost = 100

	// SlowlorisWriteInterval is the interval between partial header writes in SLOW attack
	SlowlorisWriteInterval = 10 * time.Second
)

// Constants for Proxy checking
const (
	// MinProxyCheckThreads is the minimum number of threads for proxy checking
	MinProxyCheckThreads = 100

	// MaxProxyCheckThreads is the maximum number of threads for proxy checking
	MaxProxyCheckThreads = 500

	// DefaultProxyCheckTimeout is the default timeout for checking proxy connectivity
	DefaultProxyCheckTimeout = 10 * time.Second
)
