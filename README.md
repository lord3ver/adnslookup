ADnsLookup is a DNS lookup tool written in Go language that allows to extract ip addresses.

Concurrent with worker, pool of workers and a load balance mechanism.

## Install
```sh
go get github.com/lord3ver/adnslookup/cmd/adnslookup
```

### Go >= 1.17
Installing executables with `go get` is deprecated. `go install` may be used instead:
```sh
go install github.com/lord3ver/adnslookup/cmd/adnslookup@latest
```

## Usage
```sh
adnslookup -h
```
Shows help for using the tool.
```sh

     ___    ____             __                __
    /   |  / __ \____  _____/ /   ____  ____  / /____  ______
   / /| | / / / / __ \/ ___/ /   / __ \/ __ \/ //_/ / / / __ \
  / ___ |/ /_/ / / / (__  ) /___/ /_/ / /_/ / ,< / /_/ / /_/ /
 /_/  |_/_____/_/ /_/____/_____/\____/\____/_/|_|\__,_/ .___/
                                                     /_/

A DNS Lookup tool. Concurrent with worker, pool of workers and generator "load balance" mechanism.

Version:        1.0.0
Author:         LordEver (@Lord3ver)

Usage of ADnsLookup:
  -d string
        Target domain
  -f string
        Targets file. One per line.
  -out
        Print results to stdout (default true)
  -outfile string
        Specify an output file when completed. Create or append if exists.
  -outfileNoDNS string
        Specify an output file for domains with no DNS record. Create or append if exists.
  -t int
        Max threads (default 25)
```

## Example
```sh
adnslookup -d microsoft.com

     ___    ____             __                __
    /   |  / __ \____  _____/ /   ____  ____  / /____  ______
   / /| | / / / / __ \/ ___/ /   / __ \/ __ \/ //_/ / / / __ \
  / ___ |/ /_/ / / / (__  ) /___/ /_/ / /_/ / ,< / /_/ / /_/ /
 /_/  |_/_____/_/ /_/____/_____/\____/\____/_/|_|\__,_/ .___/
                                                     /_/

A DNS Lookup tool. Concurrent with worker, pool of workers and generator "load balance" mechanism.

Version:        1.0.0
Author:         LordEver (@Lord3ver)

40.76.4.15      microsoft.com
104.215.148.63  microsoft.com
13.77.161.179   microsoft.com
40.113.200.201  microsoft.com
40.112.72.205   microsoft.com

Done!
```

## ADnsLookup Go library
ADnsLookup can be used as a library. Here is an example:

```go
package main

import (
	"fmt"
	"github.com/lord3ver/adnslookup/pkg/dnsip"
)

func main() {
	targets := []string{"www.google.com", "www.microsoft.com"}

	res, resNone := dnsip.Lookup(1, targets, false)

	fmt.Print("Found:\n")
	for _, found := range res {
		fmt.Println(found)
	}

	fmt.Print("\nNot Found:\n")
	for _, notFound := range resNone {
		fmt.Println(notFound)
	}
}
```