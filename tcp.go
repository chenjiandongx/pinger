package pinger

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// TCPPingOpts is the option set for the TCP Ping.
type TCPPingOpts struct {
	// PingTimeout is the timeout for a ping request.
	PingTimeout time.Duration
	// PingCount is the number of requests that will be sent to compute the ping quality of a host.
	PingCount int
	// MaxConcurrency sets the maximum goroutine used.
	MaxConcurrency int
	// Interval returns a time.Duration as the delay.
	Interval func() time.Duration
}

// DefaultTCPPingOpts will be used if PingOpts is nil with the TCPPing function.
var DefaultTCPPingOpts = &TCPPingOpts{
	PingTimeout:    3 * time.Second,
	PingCount:      10,
	Interval:       defaultInterval,
	MaxConcurrency: 10,
}

func (opts *TCPPingOpts) ping(dest *destination, args ...interface{}) {
	now := time.Now()
	conn, err := net.DialTimeout("tcp", dest.host, opts.PingTimeout)
	if err != nil {
		dest.addResult(zeroDur, err)
		logrus.Warnf("ping host(%s) error: %+v", dest.host, err)
		return
	}
	defer conn.Close()
	dest.addResult(time.Since(now), nil)
}

func TCPPing(opts *TCPPingOpts, hosts ...string) ([]PingStat, error) {
	if opts == nil {
		opts = DefaultTCPPingOpts
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
		pingCount:      opts.PingCount,
		ping:           opts.ping,
		setInterval:    opts.Interval,
		dest:           dests,
	})

	return sortHosts(stats, hosts...), nil
}
