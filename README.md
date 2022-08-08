# cookiestxt

[![Build Status](https://travis-ci.org/mengzhuo/cookiestxt.svg?branch=master)](https://travis-ci.org/mengzhuo/cookiestxt)
[![Go Reference](https://pkg.go.dev/badge/github.com/mengzhuo/cookietxt.svg)](https://pkg.go.dev/github.com/mengzhuo/cookietxt)
[![GoReportCard](https://goreportcard.com/badge/github.com/mengzhuo/cookiestxt)](https://goreportcard.com/report/github.com/mengzhuo/cookiestxt)
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
