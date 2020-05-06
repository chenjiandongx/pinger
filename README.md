# Pinger

> 📌 A portable ping library written in Go.

Pinger is a library used to evaluate the quality of services with ICMP/TCP/HTTP protocol.

[![GoDoc](https://godoc.org/github.com/chenjiandongx/pinger?status.svg)](https://godoc.org/github.com/chenjiandongx/pinger)
[![Travis](https://travis-ci.org/go-echarts/go-echarts.svg?branch=master)](https://travis-ci.org/chenjiandongx/pinger)
[![Appveyor](https://ci.appveyor.com/api/projects/status/v7w3u0p66grbfpxb/branch/master?svg=true)](https://ci.appveyor.com/project/chenjiandongx/pinger/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenjiandongx/pinger)](https://goreportcard.com/report/github.com/chenjiandongx/pinger)
[![License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

### 🔰 Installation

```shell
$ go get -u github.com/chenjiangdongx/pinger
```

### 📝 Usage

> For more information, please refer to [godoc](https://godoc.org/github.com/chenjiandongx/pinger).

```golang
package main

import (
	"fmt"

	"github.com/chenjiandongx/pinger"
)

func main() {
	// ICMP
	stats, err := pinger.ICMPPing(nil, "qq.com", "baidu.com", "114.114.114.114")

	// TCP
	// stats, err := pinger.TCPPing(nil, "baidu.com:80", "qq.com:80", "qq.com:443", "baidu.com:443")

	// HTTP/HTTPS
	// stats, err := pinger.HTTPPing(nil, "http://baidu.com", "https://baidu.com", "http://39.156.69.79")
	if err != nil {
		panic(err)
	}

	for _, s := range stats {
		fmt.Printf("%+v\n", s)
	}
}

// output
{Host:qq.com PktSent:10 PktLoss:0 Mean:42.858058ms Last:39.66054ms Best:38.031497ms Worst:71.050511ms}
{Host:baidu.com PktSent:10 PktLoss:0 Mean:45.834938ms Last:44.408987ms Best:40.155878ms Worst:75.480914ms}
{Host:114.114.114.114 PktSent:10 PktLoss:0 Mean:10.953486ms Last:6.618554ms Best:5.407619ms Worst:38.53662ms}
```

### 📃 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
