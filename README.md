# Pinger

> A portable ping library written in Go.

Pinger is a library used to evaluate the quality of services with ICMP/TCP/HTTP protocol.

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
	// IMMP
	stats, err := pinger.ICMPPing(nil, "qq.com", "baidu.com", "114.114.114.114")

	// TCP
	// stats, err := pinger.TCPPing(nil, "baidu.com:80", "qq.com:80", "qq.com:443", "baidu.com:443")

	// HTTP
	// stats, err := pinger.HTTPPing(nil, "baidu.com", "qq.com", "huya.com")
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
