package pinger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidICMPPing(t *testing.T) {
	hosts := []string{"qq.com", "baidu.com", "114.114.114.114"}
	stats, err := ICMPPing(nil, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, len(hosts))
}

func TestInvalidICMPPing(t *testing.T) {
	hosts := []string{"114.114.114.115"}
	opts := DefaultICMPPingOpts
	opts.PingCount = 2
	opts.PingTimeout = 200 * time.Millisecond
	stats, err := ICMPPing(opts, hosts...)
	assert.NoError(t, err)
	assert.Len(t, stats, 1)
	assert.Equal(t, stats[0].PktLoss, float64(1))
}
