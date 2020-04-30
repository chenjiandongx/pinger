package pinger

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/digineo/go-ping"
	"github.com/sirupsen/logrus"
)

type PingOpts struct {
	PingTimeout     time.Duration
	ResolverTimeout time.Duration
	Bind4           string
	Interval        func() time.Duration
	PayloadSize     uint16
	StatBufferSize  int
	Count           int
}

var DefaultPingOpts = &PingOpts{
	PingTimeout:     3 * time.Second,
	Bind4:           "0.0.0.0",
	ResolverTimeout: 1500 * time.Millisecond,
	Interval:        func() time.Duration { return time.Duration(rand.Int63n(300)) * time.Millisecond },
	PayloadSize:     56,
	StatBufferSize:  50,
	Count:           10,
}

type PingStat struct {
	PktSent int
	PktLoss float64
	Mean    time.Duration
	Last    time.Duration
	Best    time.Duration
	Worst   time.Duration
}

func Ping(opts *PingOpts, hosts ...string) (map[string]PingStat, error) {
	if opts == nil {
		opts = DefaultPingOpts
	}

	stats := make(map[string]PingStat)

	pinger := &ping.Pinger{}
	instance, err := ping.New(opts.Bind4, "::")
	if err != nil {
		return stats, fmt.Errorf("init pinger error: %s", err.Error())
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
			logrus.Warnf("resolve address error:+%v\n", err)
			continue
		}

		for _, remote := range remotes {
			// only ipv4
			if remote.IP.To4() == nil {
				continue
			}

			ipaddr := remote // need to create a copy
			dst := destination{
				host:    host,
				remote:  &ipaddr,
				history: &history{results: make([]time.Duration, opts.StatBufferSize)},
			}
			dests = append(dests, &dst)
		}
	}

	mux := sync.Mutex{}
	wg := sync.WaitGroup{}
	for c := 0; c < opts.Count; c++ {
		for _, d := range dests {
			wg.Add(1)
			go func(d *destination) {
				defer wg.Done()
				d.ping(pinger, opts.PingTimeout)

				mux.Lock()
				stats[d.host] = d.compute()
				mux.Unlock()
			}(d)
		}
		wg.Wait()
		time.Sleep(opts.Interval())
	}

	return stats, nil
}

type destination struct {
	host    string
	remote  *net.IPAddr
	display string
	*history
}

func (u *destination) ping(pinger *ping.Pinger, timeout time.Duration) {
	rtt, err := pinger.Ping(u.remote, timeout)
	if err != nil {
		logrus.Warnf("ping host[%s] error: %+v", u.host, err)
	}
	u.addResult(rtt, err)
}

type history struct {
	received int
	lost     int
	results  []time.Duration // ring, start index = .received%len
	mtx      sync.RWMutex
}

func (s *history) addResult(rtt time.Duration, err error) {
	s.mtx.Lock()
	if err == nil {
		s.results[s.received%len(s.results)] = rtt
		s.received++
	} else {
		s.lost++
	}
	s.mtx.Unlock()
}

func (s *history) compute() (st PingStat) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.received == 0 {
		if s.lost > 0 {
			st.PktLoss = 1.0
		}
		return
	}

	collection := s.results[:]
	st.PktSent = s.received + s.lost
	size := len(s.results)
	st.Last = collection[(s.received-1)%size]

	// we don't yet have filled the buffer
	if s.received <= size {
		collection = s.results[:s.received]
		size = s.received
	}

	st.PktLoss = float64(s.lost) / float64(s.received+s.lost)
	st.Best, st.Worst = collection[0], collection[0]

	total := time.Duration(0)
	for _, rtt := range collection {
		if rtt < st.Best {
			st.Best = rtt
		}
		if rtt > st.Worst {
			st.Worst = rtt
		}
		total += rtt
	}

	st.Mean = time.Duration(float64(total) / float64(size))
	return
}

func resolve(addr string, timeout time.Duration) ([]net.IPAddr, error) {
	if strings.ContainsRune(addr, '%') {
		ipaddr, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			return nil, err
		}
		return []net.IPAddr{*ipaddr}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return net.DefaultResolver.LookupIPAddr(ctx, addr)
}
