# Pinger

```golang
package main

import (
	"fmt"

	"github.com/chenjiandongx/pinger"
)

func main() {
	stats, err := pinger.Ping(nil, "huya.com", "qq.com", "114.114.114.114")
	if err != nil {
		panic(err)
	}

	for k, v := range stats {
		fmt.Printf("[host]:%s; [stats]:%+v\n", k, v)
	}
}
```
