# cookiestxt

[![Build Status](https://github.com/mengzhuo/cookiestxt/actions/workflows/go.yml/badge.svg)](https://github.com/mengzhuo/cookiestxt/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/mengzhuo/cookiestxt.svg)](https://pkg.go.dev/github.com/mengzhuo/cookiestxt)
[![GoReportCard](https://goreportcard.com/badge/github.com/mengzhuo/cookiestxt)](https://goreportcard.com/report/github.com/mengzhuo/cookiestxt)
[![Coverage Status](https://coveralls.io/repos/github/mengzhuo/cookiestxt/badge.svg?branch=master)](https://coveralls.io/github/mengzhuo/cookiestxt?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

cookiestxt implement parser of cookies txt format

## Usage

```golang
package main

import (
        "log"
        "strings"

        "github.com/mengzhuo/cookiestxt"
)

func main() {
        buf := strings.NewReader(`
        # This example taken from http://www.cookiecentral.com/faq/#3.5
        #HttpOnly_.netscape.com TRUE / FALSE 946684799 NETSCAPE_ID 100103
        `)
        cl, err := cookiestxt.Parse(buf)
        log.Print(cl, err)
}
```

```
$ go run main.go 

[NETSCAPE_ID=100103; Path=/; Domain=netscape.com; Expires=Fri, 31 Dec 1999 23:59:59
 GMT; HttpOnly] <nil>
```
