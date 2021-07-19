# Pinger

> 📌 A portable ping library written in Go.

Pinger is a library used to evaluate the quality of services in ICMP/TCP/HTTP protocol.

What's worth pointing out is that `ping` here means not only the standard **IMCP Protocol**, but also **TCP/HTTP/HTTPS**. It's the more general sense here as an approach to detect the network quality.

[![GoDoc](https://godoc.org/github.com/chenjiandongx/pinger?status.svg)](https://godoc.org/github.com/chenjiandongx/pinger)
[![Travis](https://travis-ci.org/chenjiandongx/pinger.svg?branch=master)](https://travis-ci.org/chenjiandongx/pinger)
[![Appveyor](https://ci.appveyor.com/api/projects/status/v7w3u0p66grbfpxb/branch/master?svg=true)](https://ci.appveyor.com/project/chenjiandongx/pinger/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenjiandongx/pinger)](https://goreportcard.com/report/github.com/chenjiandongx/pinger)
[![License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

### 🔰 Installation

```shell
$ go get -u github.com/chenjiandongx/pinger
```

### 📝 Basic Usage

```golang
package main

import (
	"fmt"

	"github.com/chenjiandongx/pinger"
)

func main() {
	// ICMP
	stats, err := pinger.ICMPPing(nil, "huya.com", "youtube.com", "114.114.114.114")

	// TCP
	// stats, err := pinger.TCPPing(nil, "huya.com:80", "google.com:80", "huya.com:443", "google.com:443")

	// HTTP/HTTPS
	// stats, err := pinger.HTTPPing(nil, "http://huya.com", "https://google.com", "http://39.156.69.79")
	if err != nil {
		panic(err)
	}

	for _, s := range stats {
		fmt.Printf("%+v\n", s)
	}
}

// output
{Host:huya.com PktSent:10 PktLossRate:0 Mean:36.217604ms Last:37.636144ms Best:34.949608ms Worst:37.636144ms}
{Host:youtube.com PktSent:10 PktLossRate:0 Mean:12.853867ms Last:13.105932ms Best:11.836598ms Worst:14.844776ms}
{Host:114.114.114.114 PktSent:10 PktLossRate:0 Mean:7.689701ms Last:7.711122ms Best:6.252377ms Worst:9.427482ms}
```

### 🎉 Options

**PingTimeout / Interval**
```golang
opts := pinger.DefaultICMPPingOpts
opts.PingTimeout = 50 * time.Millisecond
// sleep 60 mills after a request every time.
opts.Interval = func() time.Duration { return 60 * time.Millisecond }
```

**PingCount / MaxConcurrency / FailOver**
```golang
opts := pinger.DefaultICMPPingOpts
// network is unstable thus we need more ping-ops to evaluate the network quality overall.
opts.PingCount = 20
// set the maximum concurreny, goroutine is cheap, but not free :)
opts.MaxConcurrency = 5
// set per host maximum failed allowed
opts.FailOver = 5
```

**Headers / Method**
```golang
// in http/https case, there are more options could be used.
opts := pinger.DefaultHTTPPingOpts
// HTTP headers, something speical for authentication or anything else.
opts.Headers = map[string]string{"token": "my-token", "who": "me"}
// HTTP Method, feel free to use any standard HTTP methods you need.
opts.Method = http.MethodPost
```

*For more information, please refer to [the documentation](https://godoc.org/github.com/chenjiandongx/pinger).*

### 📃 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
