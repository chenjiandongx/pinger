package pinger

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPPingOpts is the option set for the HTTP Ping.
type HTTPPingOpts struct {
	// PingTimeout is the timeout for a ping request.
	PingTimeout time.Duration
	// PingCount is the number of requests that will be sent to compute the ping quality of a host.
	PingCount int
	// MaxConcurrency sets the maximum goroutine used.
	MaxConcurrency int
	// Interval returns a time.Duration as the delay.
	Interval func() time.Duration

	// Method represents the HTTP Method(GET/POST/PUT/...).
	Method string
	// Body represents the HTTP Request body.
	Body io.Reader
	// Headers represents for the HTTP Headers.
	Headers map[string]string
}

// DefaultHTTPPingOpts will be used if PingOpts is nil with the HTTPPing function.
var DefaultHTTPPingOpts = &HTTPPingOpts{
	PingTimeout:    3 * time.Second,
	PingCount:      10,
	Method:         http.MethodGet,
	Body:           nil,
	Headers:        nil,
	Interval:       func() time.Duration { return time.Duration(rand.Int63n(200)) * time.Millisecond },
	MaxConcurrency: 10,
}

func (opts *HTTPPingOpts) ping(dest *destination, args ...interface{}) {
	client := args[0].(*http.Client)
	check := func(err error) {
		logrus.Warnf("ping host(%s) error: %+v", dest.host, err)
		dest.addResult(zeroDur, err)
	}

	req, err := http.NewRequest(opts.Method, dest.host, opts.Body)
	if err != nil {
		check(err)
		return
	}

	if opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Add(k, v)
		}
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	now := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		check(err)
		return
	}

	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
		check(err)
		return
	}
	defer resp.Body.Close()
	dest.addResult(time.Since(now), nil)
}

func HTTPPing(opts *HTTPPingOpts, hosts ...string) ([]PingStat, error) {
	if opts == nil {
		opts = DefaultHTTPPingOpts
	}

	var transport = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext:       (&net.Dialer{Timeout: opts.PingTimeout, KeepAlive: time.Second}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   opts.PingTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
		args:           client,
	})

	return sortHosts(stats, hosts...), nil
}
