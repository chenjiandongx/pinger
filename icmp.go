package pinger

import (
	"fmt"
	"time"

	"github.com/digineo/go-ping"
	"github.com/sirupsen/logrus"
)

// ICMPPingOpts is the option set for the ICMP Ping.
type ICMPPingOpts struct {
	// PingTimeout is the timeout for a ping request.
	PingTimeout time.Duration
	// PingCount is the number of requests that will be sent to compute the ping quality of a host.
	PingCount int
	// MaxConcurrency sets the maximum goroutine used.
	MaxConcurrency int
	// Interval returns a time.Duration as the delay.
	Interval func() time.Duration

	// ResolverTimeout is the timeout for the net.ResolveIPAddr request.
	ResolverTimeout time.Duration
	// Bind4 is the ipv4 bind to start a raw socket.
	Bind4 string
	// PayloadSize represents the request body size for a ping request.
	PayloadSize uint16
}

// DefaultICMPPingOpts will be used if PingOpts is nil with the ICMPPing function.
var DefaultICMPPingOpts = &ICMPPingOpts{
	PingTimeout:     3 * time.Second,
	PingCount:       10,
	MaxConcurrency:  10,
	Interval:        defaultInterval,
	Bind4:           "0.0.0.0",
	ResolverTimeout: 1500 * time.Millisecond,
	PayloadSize:     56,
}

func (opts *ICMPPingOpts) ping(dest *destination, args ...interface{}) {
	pinger := args[0].(*ping.Pinger)
	rtt, err := pinger.Ping(dest.remote, opts.PingTimeout)
	if err != nil {
		logrus.Warnf("ping host(%s) error: %+v", dest.host, err)
	}
	dest.addResult(rtt, err)
}

func ICMPPing(opts *ICMPPingOpts, hosts ...string) ([]PingStat, error) {
	if opts == nil {
		opts = DefaultICMPPingOpts
	}

	pinger := &ping.Pinger{}
	instance, err := ping.New(opts.Bind4, "::")
	if err != nil {
		return nil, fmt.Errorf("init pinger error: %s", err.Error())
	}

	if instance.PayloadSize() != opts.PayloadSize {
		instance.SetPayloadSize(opts.PayloadSize)
	}
	pinger = instance
	defer pinger.Close()

	dests := make([]*destination, 0)
	for _, host := range hosts {
		remotes, err := resolve(host, opts.ResolverTimeout)
		if err != nil {
			logrus.Warnf("resolve address error: %+v\n", err)
			continue
		}

		for _, remote := range remotes {
			// only ipv4
			if remote.IP.To4() == nil {
				continue
			}

			ipaddr := remote // need to create a copy
			dests = append(dests, &destination{
				host:    host,
				remote:  &ipaddr,
				history: &history{results: make([]time.Duration, defaultStatsBuf)},
			})
		}
	}

	stats := calculateStats(calcStatsReq{
		maxConcurrency: opts.MaxConcurrency,
		pingCount:      opts.PingCount,
		ping:           opts.ping,
		setInterval:    opts.Interval,
		dest:           dests,
		args:           pinger,
	})

	return sortHosts(stats, hosts...), nil
}
