package curl

import (
	"fmt"
	"net/http/httputil"
	"slices"
	"testing"
)

func TestCurlToRequest(t *testing.T) {

	test := `curl 'https://waa-pa.clients6.google.com/$rpc/google.internal.waa.v1.Waa/Create' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'X-User-Agent: grpc-web-javascript/0.1' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'Authorization: SAPISIDHASH 1712976476_a588e7ff1e4623decb92fd59c86ea4700e475013' \
	-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36' \
	-H 'Content-Type: application/json+protobuf' \
	-H 'Referer: https://www.google.com/' \
	-H 'X-Goog-Api-Key: AIzaSyBGb5fGAyC-pRcRU6MUHb__b_vKha71HRE' \
	-H 'X-Goog-AuthUser: 0' \
	-H 'sec-ch-ua-platform: "macOS"' \
	--data-raw '["/JR8jsAkqotcKsEKhXic"]'`

	req, err := CurlToRequest(test)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}

func TestCurlToRequestMac(t *testing.T) {
	test := `curl -H "Host: connectivitycheck.gstatic.com" -H "Accept: */*" -H "User-Agent: Twitter/19.9.30 CFNetwork/1492.0.1 Darwin/23.3.0" -H "Accept-Language: zh-CN,zh-Hans;q=0.9" --compressed "http://connectivitycheck.gstatic.com/generate_204"`

	req, err := CurlToRequest(test)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))

}

func TestParseArgs(t *testing.T) {
	t.Run("base1", func(t *testing.T) {
		x := parseArgs("curl -H 'are you ok'")
		if !slices.Equal(x, []string{"curl", "-H", "are you ok"}) {
			t.Errorf("test fail %v", x)
		}
	})
	t.Run("base2", func(t *testing.T) {
		x := parseArgs(`curl -H "are you ok"`)
		if !slices.Equal(x, []string{"curl", "-H", "are you ok"}) {
			t.Errorf("test fail %v", x)
		}
	})

	t.Run("base3", func(t *testing.T) {
		x := parseArgs(`curl -H 'are you="hello" ok'`)
		if !slices.Equal(x, []string{"curl", "-H", `are you="hello" ok`}) {
			t.Errorf("test fail %v", x)
		}
	})

	t.Run("base4", func(t *testing.T) {
		x := parseArgs(`curl -H 'are you= "hello" ok'`)
		if !slices.Equal(x, []string{"curl", "-H", `are you= "hello" ok`}) {
			t.Errorf("test fail %v", x)
		}
	})

	t.Run("base5", func(t *testing.T) {
		x := parseArgs(`curl -H 'are you= "hello ok"'`)
		if !slices.Equal(x, []string{"curl", "-H", `are you= "hello ok"`}) {
			t.Errorf("test fail %v", x)
		}
	})

}
