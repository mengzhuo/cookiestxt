// Copyright 2017 Meng Zhuo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cookiestxt implement parser of cookies txt format that commonly supported by
// curl / wget / aria2c / chrome / firefox
//
// see http://www.cookiecentral.com/faq/#3.5 for more detail
package cookiestxt

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// http://www.cookiecentral.com/faq/#3.5
	// The domain that created AND that can read the variable.
	domainIdx = iota
	// A TRUE/FALSE value indicating if all machines within a given domain can access the variable. This value is set automatically by the browser, depending on the value you set for domain.
	flagIdx
	// The path within the domain that the variable is valid for.
	pathIdx
	// A TRUE/FALSE value indicating if a secure connection with the domain is needed to access the variable.
	secureIdx
	// The UNIX time that the variable will expire on. UNIX time is defined as the number of seconds since Jan 1, 1970 00:00:00 GMT.
	expirationIdx
	// The name of the variable.
	nameIdx
	// The value of the variable.
	valueIdx
)

const (
	httpOnlyPrefix = "#HttpOnly_"
	fieldsCount    = 7
)

// Parse cookie txt file format from input stream
func Parse(rd io.Reader) (cl []*http.Cookie, err error) {
	scanner := bufio.NewScanner(rd)
	// allow bigger lines if needed
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var line int
	for scanner.Scan() {
		line++

		trimed := strings.TrimSpace(scanner.Text())
		if trimed == "" {
			continue
		}

		// skip comments except HttpOnly_ prefixed lines
		if strings.HasPrefix(trimed, "#") && !strings.HasPrefix(trimed, httpOnlyPrefix) {
			continue
		}

		var c *http.Cookie
		c, err = ParseLine(trimed)
		if err != nil {
			return cl, fmt.Errorf("cookiestxt line:%d, err:%s", line, err)
		}
		cl = append(cl, c)
	}

	err = scanner.Err()
	return
}

// ParseLine parse single cookie from one line with stricter validation
func ParseLine(raw string) (c *http.Cookie, err error) {
	raw = strings.TrimSpace(raw)
	f := strings.Fields(raw)
	if len(f) == fieldsCount-1 {
		// missing value -> treat as empty
		f = append(f, "")
	} else if len(f) < fieldsCount {
		err = fmt.Errorf("expecting fields=%d, got=%d", fieldsCount, len(f))
		return
	}

	// basic required fields
	if strings.TrimSpace(f[domainIdx]) == "" {
		err = fmt.Errorf("empty domain")
		return
	}
	if strings.TrimSpace(f[nameIdx]) == "" {
		err = fmt.Errorf("empty cookie name")
		return
	}

	// validate flag (second field) format but do not use value
	if _, perr := parseBoolStrict(f[flagIdx]); perr != nil {
		err = fmt.Errorf("invalid flag value: %v", perr)
		return
	}

	// secure field must be a valid boolean token
	secureVal, perr := parseBoolStrict(f[secureIdx])
	if perr != nil {
		err = fmt.Errorf("invalid secure value: %v", perr)
		return
	}

	c = &http.Cookie{
		Raw:    raw,
		Name:   f[nameIdx],
		Value:  f[valueIdx],
		Path:   f[pathIdx],
		MaxAge: 0,
		Secure: secureVal,
	}

	var ts int64
	ts, err = strconv.ParseInt(f[expirationIdx], 10, 64)
	if err != nil {
		return
	}
	c.Expires = time.Unix(ts, 0)

	c.Domain = f[domainIdx]
	if strings.HasPrefix(c.Domain, httpOnlyPrefix) {
		c.HttpOnly = true
		c.Domain = c.Domain[len(httpOnlyPrefix):]
	}

	return
}

// parseBoolStrict validates boolean tokens and returns an error on unknown token.
// Accepts: "1"/"0", "TRUE"/"FALSE" (case-insensitive).
func parseBoolStrict(input string) (bool, error) {
	s := strings.TrimSpace(input)
	if s == "1" || s == "0" {
		return s == "1", nil
	}
	if strings.EqualFold(s, "TRUE") {
		return true, nil
	}
	if strings.EqualFold(s, "FALSE") {
		return false, nil
	}
	return false, fmt.Errorf("expect TRUE/FALSE or 1/0, got %q", input)
}

// parseBool kept for compatibility; returns false on invalid input
func parseBool(input string) bool {
	b, _ := parseBoolStrict(input)
	return b
}
