# Pinger

```golang
package main

import (
	"fmt"

	"github.com/chenjiandongx/pinger"
)

func main() {
	stats, err := pinger.HTTPPing(nil, "baidu.com", "qq.com", "huya.com")
	if err != nil {
		panic(err)
	}

	for _, s := range stats {
		fmt.Printf("%+v\n", s)
	}
}

// output
{Host:baidu.com PktSent:10 PktLoss:0 Mean:59.139136ms Last:52.84319ms Best:49.859695ms Worst:84.807026ms}
{Host:qq.com PktSent:10 PktLoss:0 Mean:52.69315ms Last:48.398119ms Best:46.002873ms Worst:64.527586ms}
{Host:huya.com PktSent:10 PktLoss:0 Mean:26.389519ms Last:27.576105ms Best:19.517054ms Worst:39.46555ms}
```
