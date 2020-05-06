package pinger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidTCPPing(t *testing.T) {
	hosts := []string{"baidu.com:80", "qq.com:80", "qq.com:443", "baidu.com:443"}
	stats, err := TCPPing(nil, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, len(hosts))
}

func TestInvalidTCPPing(t *testing.T) {
	hosts := []string{"114.114.115.115:8080"}
	opts := DefaultTCPPingOpts
	opts.PingCount = 2
	opts.PingTimeout = 200 * time.Millisecond
	stats, err := TCPPing(opts, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, 1)
	assert.Equal(t, stats[0].PktLossRate, float64(1))
}
