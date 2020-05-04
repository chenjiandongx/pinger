package pinger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidHTTPPing(t *testing.T) {
	hosts := []string{"http://baidu.com", "https://baidu.com", "http://39.156.69.79"}
	stats, err := HTTPPing(nil, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, len(hosts))
}

func TestInvalidHTTPPing(t *testing.T) {
	hosts := []string{"http://114.114.115.115"}
	opts := DefaultHTTPPingOpts
	opts.PingCount = 2
	opts.PingTimeout = 200 * time.Millisecond
	stats, err := HTTPPing(opts, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, 1)
	assert.Equal(t, stats[0].PktLoss, float64(1))
}
