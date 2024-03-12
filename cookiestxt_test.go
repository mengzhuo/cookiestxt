package cookiestxt

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func ExampleParseLine() {

	c, _ := ParseLine(".netscape.com TRUE / TRUE 946684799 NETSCAPE_ID 100105")
	fmt.Printf("Name=%s, Value=%s\n", c.Name, c.Value)
	// Output: Name=NETSCAPE_ID, Value=100105
}

func ExampleParse() {

	buf := strings.NewReader(`
        # This example taken from http://www.cookiecentral.com/faq/#3.5
        #HttpOnly_.netscape.com TRUE / FALSE 946684799 NETSCAPE_ID 100105
        `)
	cl, _ := Parse(buf)
	fmt.Printf("Name[0]=%s, Value[0]=%s, Len=%d", cl[0].Name, cl[0].Value, len(cl))
	// Output: Name[0]=NETSCAPE_ID, Value[0]=100105, Len=1
}

func TestParseEmptyValue(t *testing.T) {
	bt := time.Unix(946684799, 0)
	c, err := ParseLine(".netscape.com TRUE / TRUE 946684799 NETSCAPE_ID ")
	if err != nil || c.Value != "" || c.Expires.Before(bt) || !c.Secure {
		t.Error(err)
	}
	t.Log(c)
}

func TestParseLine(t *testing.T) {
	bt := time.Unix(946684799, 0)
	c, err := ParseLine(".netscape.com TRUE / TRUE 946684799 NETSCAPE_ID 100103")
	if err != nil || c.Value != "100103" || c.Expires.Before(bt) || !c.Secure {
		t.Error(err)
	}
	t.Log(c)
	_, err = ParseLine(".netscape.com / FALSE 946684799 NETSCAPE_ID 100103")
	if err == nil {
		t.Error("no error on invalid txt")
	}
}

func TestParseHTTPOnlyLine(t *testing.T) {
	c, err := ParseLine("#HttpOnly_.netscape.com TRUE / FALSE 946684799 NETSCAPE_ID 100103")
	if err != nil || c.Domain != ".netscape.com" || c.HttpOnly == false {
		t.Error(err)
	}
	t.Log(c)
}

func TestParseUnixTime(t *testing.T) {
	c, err := ParseLine("#HttpOnly_.netscape.com TRUE / FALSE -1 NETSCAPE_ID 100103")
	if err != nil || c.Expires.After(time.Unix(0, 0)) {
		t.Error(err)
	}
	t.Log(c)

	_, err = ParseLine("#HttpOnly_.netscape.com TRUE / FALSE NOT_A_INT NETSCAPE_ID 100103")
	if err == nil {
		t.Error("no error on invalid expires")
	}
}

func TestParseFunc(t *testing.T) {
	mock := `
	# Comment
	# More Comment
	# This is a long comment that_has_7 fields

	#not very well comment
	#HttpOnly_.netscape.com TRUE / FALSE 946684799 NETSCAPE_ID 100103
	.netscape.com TRUE / FALSE 946684799 NETSCAPE_ID 100103
	#

	`
	cl, err := Parse(strings.NewReader(mock))
	if err != nil || len(cl) != 2 {
		t.Error(err, cl)
	}
}

func TestParseFailed(t *testing.T) {
	mock := `
	#
	#HttpOnly_.netscape.com TRUE / FALSE NOT_A_INT NETSCAPE_ID 100103
	#

	`
	_, err := Parse(strings.NewReader(mock))
	if err == nil || strings.Index(err.Error(), "line:3") == -1 {
		t.Error(err)
	}
}

func BenchmarkParseLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseLine(".netscape.com / FALSE 946684799 NETSCAPE_ID 100103")
	}
}
