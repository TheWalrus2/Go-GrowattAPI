package growatt

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func Login(user string, passwd string) string {

	h := md5.New()
	io.WriteString(h, "YourPassword")

	data := url.Values{
		"userName": {"YourUsername"},
		"password": {hex.EncodeToString(h.Sum(nil))},
	}

	response, err := http.PostForm("https://server-api.growatt.com/newLoginAPI.do", data)

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	message := string(responseData)
	return message
}

func gettlxdetailurlcall(tlx_id string) string {
	req, err := http.NewRequest("GET", "https://server-api.growatt.com/newTlxApi.do", nil)
	if err != nil {
		log.Print(err)
	}

	q := req.URL.Query()
	q.Add("op", "getTlxDetailData")
	q.Add("id", tlx_id)
	req.URL.RawQuery = q.Encode()

	return req.URL.String()
}

func getloginurl() url.Values {
	h := md5.New()
	io.WriteString(h, "YourPassword")
	data := url.Values{
		"userName": {"YourUsername"},
		"password": {hex.EncodeToString(h.Sum(nil))},
	}
	return data
}

func Tlxdetail(tlx_id string) string {

	jar := NewJar()
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       0,
	}
	logindata := getloginurl()
	resp, _ := client.PostForm("https://server-api.growatt.com/newLoginAPI.do", logindata)
	resp.Body.Close()

	url := gettlxdetailurlcall(tlx_id)
	response, err := client.Get(url)
	responseData, err := ioutil.ReadAll(response.Body)

	response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	message := string(responseData)
	return message
}
