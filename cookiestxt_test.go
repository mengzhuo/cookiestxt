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
		_, err := ParseLine(".netscape.com / FALSE 946684799 NETSCAPE_ID 100103")
		if err != nil {
			b.Error(err)
		}
	}
}

// ---- New stricter validation tests ----

func TestParseBoolVariants(t *testing.T) {
	cases := []struct {
		secure string
		want   bool
	}{
		{"1", true},
		{"0", false},
		{"TRUE", true},
		{"false", false},
		{"TrUe", true},
	}

	for _, tc := range cases {
		line := fmt.Sprintf(".netscape.com TRUE / %s 946684799 NETSCAPE_ID 100103", tc.secure)
		c, err := ParseLine(line)
		if err != nil {
			t.Fatalf("unexpected error for secure=%s: %v", tc.secure, err)
		}
		if c.Secure != tc.want {
			t.Fatalf("secure mismatch for %s: want %v got %v", tc.secure, tc.want, c.Secure)
		}
	}
}

func TestParseInvalidFlag(t *testing.T) {
	// invalid flag token (flagIdx)
	_, err := ParseLine(".netscape.com MAYBE / TRUE 946684799 NETSCAPE_ID 100103")
	if err == nil {
		t.Fatal("expected error on invalid flag token")
	}
}

func TestParseInvalidSecure(t *testing.T) {
	_, err := ParseLine(".netscape.com TRUE / MAYBE 946684799 NETSCAPE_ID 100103")
	if err == nil {
		t.Fatal("expected error on invalid secure token")
	}
}

func TestParseLongValue(t *testing.T) {
	longVal := strings.Repeat("a", 200000) // 200KB value
	line := fmt.Sprintf(".netscape.com TRUE / TRUE 946684799 NETSCAPE_ID %s\n", longVal)
	cl, err := Parse(strings.NewReader(line))
	if err != nil {
		t.Fatalf("unexpected error parsing long line: %v", err)
	}
	if len(cl) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cl))
	}
	if cl[0].Value != longVal {
		t.Fatal("cookie value mismatch for long value")
	}
}

func FuzzParseBool(f *testing.F) {
	cases := []string{"0", "1", "True", "False"}
	for _, tc := range cases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, s string) {
		parseBoolStrict(s)
	})
}
