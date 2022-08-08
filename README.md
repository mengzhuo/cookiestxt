# cookiestxt

[![Build Status](https://travis-ci.org/mengzhuo/cookiestxt.svg?branch=master)](https://travis-ci.org/mengzhuo/cookiestxt)
[![GoDoc](https://godoc.org/github.com/mengzhuo/cookiestxt?status.svg)](https://godoc.org/github.com/mengzhuo/cookiestxt)
[![GoReportCard](https://goreportcard.com/badge/github.com/mengzhuo/cookiestxt)](https://goreportcard.com/report/github.com/mengzhuo/cookiestxt)

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
