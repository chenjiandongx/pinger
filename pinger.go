package pinger

import (
	"context"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	defaultInterval = func() time.Duration { return time.Duration(rand.Int63n(200)) * time.Millisecond }
	defaultStatsBuf = 60
	zeroDur         = time.Duration(0)
)

// PingStat struct is used to record the ping result.
type PingStat struct {
	Host    string
	PktSent int
	PktLoss float64
	Mean    time.Duration
	Last    time.Duration
	Best    time.Duration
	Worst   time.Duration
}

type destination struct {
	host   string
	remote *net.IPAddr
	*history
}

type history struct {
	received int
	lost     int
	results  []time.Duration // ring, start index = .received%len
	mtx      sync.RWMutex
}

func (h *history) addResult(rtt time.Duration, err error) {
	h.mtx.Lock()
	if err != nil {
		h.lost++
	} else {
		h.results[h.received%len(h.results)] = rtt
		h.received++
	}
	h.mtx.Unlock()
}

func (h *history) compute() (st PingStat) {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	if h.received == 0 {
		if h.lost > 0 {
			st.PktLoss = 1.0
		}
		return
	}

	collection := h.results[:]
	st.PktSent = h.received + h.lost
	size := len(h.results)
	st.Last = collection[(h.received-1)%size]

	// we don't yet have filled the buffer
	if h.received <= size {
		collection = h.results[:h.received]
		size = h.received
	}

	st.PktLoss = float64(h.lost) / float64(h.received+h.lost)
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

func sortHosts(stats map[string]PingStat, hosts ...string) []PingStat {
	ordered := make([]PingStat, len(hosts))
	for i, host := range hosts {
		ordered[i] = stats[host]
	}
	return ordered
}

type calcStatsReq struct {
	maxConcurrency int
	pingCount      int
	dest           []*destination
	ping           func(d *destination, args ...interface{})
	setInterval    func() time.Duration
	args           interface{}
}

func calculateStats(csr calcStatsReq) map[string]PingStat {
	stats := make(map[string]PingStat)

	mux := sync.Mutex{}
	wg := sync.WaitGroup{}
	sema := make(chan struct{}, csr.maxConcurrency)

	for c := 0; c < csr.pingCount; c++ {
		for _, dest := range csr.dest {
			sema <- struct{}{}
			wg.Add(1)
			go func(d *destination) {
				defer func() {
					wg.Done()
					<-sema
				}()

				csr.ping(d, csr.args)

				mux.Lock()
				stat := d.compute()
				stat.Host = d.host
				stats[d.host] = stat
				mux.Unlock()
			}(dest)
		}
		wg.Wait()
		time.Sleep(csr.setInterval())
	}

	return stats
}
