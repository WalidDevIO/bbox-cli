package client

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type DeviceToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type BboxClient struct {
	Client *http.Client
	Url    *url.URL
	Bearer *DeviceToken
}

func NewClient(baseUrl *url.URL) (*BboxClient, error) {
	var client http.Client
	myCookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = myCookieJar
	return &BboxClient{
		Client: &client,
		Url:    baseUrl,
	}, nil
}

func (bc *BboxClient) GetCookies() []*http.Cookie {
	return bc.Client.Jar.Cookies(bc.Url)
}

func (bc *BboxClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	u := bc.Url.JoinPath(path)
	return http.NewRequest(method, u.String(), body)
}

func (bc *BboxClient) Do(req *http.Request) (*http.Response, error) {
	return bc.Client.Do(req)
}

func (bc *BboxClient) Get(url string) (*http.Response, error) {
	return bc.Client.Get(bc.Url.JoinPath(url).String())
}

func (bc *BboxClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return bc.Client.Post(bc.Url.JoinPath(url).String(), contentType, body)
}

func (bc *BboxClient) Nat() *NatInterface {
	return &NatInterface{Client: bc}
}

func (bc *BboxClient) Firewall() *FirewallInterface {
	return &FirewallInterface{Client: bc}
}

func (bc *BboxClient) Auth() *AuthInterface {
	return &AuthInterface{Client: bc}
}
