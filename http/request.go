package http

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	k5HexMin = 1 << 16
	k5HexMax = 1 << 20
)

var requestIDFormatter = strings.NewReplacer("+", "", "-", "", ".", "")

func GenRequestID() string { return GenRequestIDWithTime(time.Now()) }

func GenRequestIDWithTime(t time.Time) string {
	n := k5HexMin + rand.Intn(k5HexMax-k5HexMin)
	return requestIDFormatter.Replace(t.Format("20060102150405.000000000-0700")) +
		fmt.Sprintf("%X", n)
}

type Request struct {
	*http.Request
}

func request(r *http.Request) *Request {
	return &Request{r}
}

// Some extended handy methods.
func (r *Request) FormValueGetter() IFormValueGetter {
	return RequestFormValueGetter(r.Request)
}

func (r *Request) IsAjax() bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (r *Request) IsUpload() bool {
	return strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data")
}

// From defacto standard HTTP header field. Format below:
// X-Forwarded-For: client, proxy1, proxy2
// [Reference](http://en.wikipedia.org/wiki/X-Forwarded-For).
func (r *Request) Proxies() []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ", ")
	}
	return []string{}
}

// Get the IP of the client.
func (r *Request) IP() string {
	ips := r.Proxies()
	if len(ips) > 0 && ips[0] != "" {
		return strings.Split(ips[0], ":")[0]
	}
	ip := strings.Split(r.RemoteAddr, ":")
	if len(ip) > 0 && ip[0] != "[" {
		return ip[0]
	}
	return "127.0.0.1"
}

// Return the accept encodings list sorted by qvalue in descending order.
func (r *Request) AcceptEncodings() []string {
	return acceptEncodings(r.Header.Get("Accept-Encoding"))
}

type vq struct {
	v string
	q float32
}

type vqDescByQ []vq

func (vq vqDescByQ) Len() int           { return len(vq) }
func (vq vqDescByQ) Swap(i, j int)      { vq[i], vq[j] = vq[j], vq[i] }
func (vq vqDescByQ) Less(i, j int) bool { return vq[i].q > vq[j].q }

var vqReg *regexp.Regexp = regexp.MustCompile("(\\w+|\\*)(;q=(0(\\.\\d{0,3})*|1(\\.0{0,3})*))*")

func extractVqs(rawstr string) (vqs []vq) {
	all := vqReg.FindAllStringSubmatch(rawstr, -1)
	for _, sub := range all {
		var q float64 = 1.0
		if sub[3] != "" {
			q, _ = strconv.ParseFloat(sub[3], 32)
		}
		vqs = append(vqs, vq{v: sub[1], q: float32(q)})
	}
	return
}

// Extract the accepted compression encodings and figure out whether `identity` is acceptable.
// See rfc-2616 14.3 Accept-Encoding.
func acceptEncodings(rawstr string) (acceptEncodings []string) {
	if rawstr == "" {
		return []string{"identity"}
	}
	vqs := extractVqs(rawstr)
	sort.Sort(vqDescByQ(vqs))

	zeroAsterisk, hasIdentity := false, false

	for _, x := range vqs {
		if x.v == "*" && x.q == 0 {
			zeroAsterisk = true
		}
		if x.v == "identity" {
			hasIdentity = true
		}
		if x.q > 0 {
			acceptEncodings = append(acceptEncodings, x.v)
		}
	}

	if !hasIdentity && !zeroAsterisk {
		acceptEncodings = append(acceptEncodings, "identity")
	}

	return
}

func (r *Request) GeneralAccessLogItems() map[string]interface{} {
	items := map[string]interface{}{
		"path":       r.URL.Path,
		"host":       r.URL.Host,
		"ip":         r.IP(),
		"method":     r.Method,
		"referer":    r.Referer(),
		"user_agent": r.UserAgent(),
	}
	values := r.URL.Query()
	for k, v := range values {
		items[k] = v
	}
	return items
}
