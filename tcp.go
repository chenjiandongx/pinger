package pinger

import (
	"math/rand"
	"net"
	"time"
)

// TCPPingOpts is the option set for the TCP Ping.
type TCPPingOpts struct {
	// PingTimeout is the timeout for a ping request.
	PingTimeout time.Duration
	// PingCount counting requests for calculating ping quality of host.
	PingCount int
	// MaxConcurrency sets the maximum goroutine used.
	MaxConcurrency int
	// FailOver is the per host maximum failed allowed.
	FailOver int
	// Interval returns a time.Duration as the delay.
	Interval func() time.Duration
}

// DefaultTCPPingOpts will be used if PingOpts is nil with the TCPPing function.
func DefaultTCPPingOpts() *TCPPingOpts {
	return &TCPPingOpts{
		PingTimeout:    3 * time.Second,
		PingCount:      10,
		Interval:       func() time.Duration { return time.Duration(rand.Int63n(200)) * time.Millisecond },
		MaxConcurrency: 10,
		FailOver:       5,
	}
}

func (opts *TCPPingOpts) ping(dest *destination, args ...interface{}) error {
	now := time.Now()
	conn, err := net.DialTimeout("tcp", dest.host, opts.PingTimeout)
	if err != nil {
		dest.addResult(zeroDur, err)
		return err
	}

	if err = conn.Close(); err != nil {
		dest.addResult(zeroDur, err)
		return err
	}
	dest.addResult(time.Since(now), nil)
	return nil
}

func TCPPing(opts *TCPPingOpts, hosts ...string) ([]PingStat, error) {
	if opts == nil {
		opts = DefaultTCPPingOpts()
	}

	dests := make([]*destination, 0)
	for _, host := range hosts {
		dests = append(dests, &destination{
			host:    host,
			remote:  nil,
			history: &history{results: make([]time.Duration, defaultStatsBuf)},
		})
	}

	stats := calculateStats(calcStatsReq{
		maxConcurrency: opts.MaxConcurrency,
		failover:       opts.FailOver,
		pingCount:      opts.PingCount,
		ping:           opts.ping,
		setInterval:    opts.Interval,
		dest:           dests,
	})

	return sortHosts(stats, hosts...), nil
}
