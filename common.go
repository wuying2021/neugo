package neugo

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

/*
	session := neu.NewSession()
	neu.Use(session).WithAuth("20185342","test").On(neu.WebVPN).Login()
	neu.Use(session).WithAuth("20185342","test").On(neu.WebVPN).LoginService("xxx")
	neu.Use(session).WithAuth("20185342","test").Validate()
	neu.Use(session).WithToken("xxx").On(neu.CAS).Login()
	neu.Use(session).WithToken("xxx").On(neu.CAS).LoginService("xxx")
	neu.Use(session).WithToken("xxx").Validate()
	neu.About(session).Token()
	neu.About(session).Info()
*/

// 获取带有cookiejar的http.Client，默认Timeout为3秒
func NewSession() *http.Client {
	n := &http.Client{Timeout: 3 * time.Second}
	jar, _ := cookiejar.New(nil)
	// 绑定session
	n.Jar = jar
	return n
}

// 平台
type Platform = byte

const (
	CAS Platform = iota
	WebVPN
)

/*
	USE
*/
// 提供登陆动作的链式调用
func Use(client *http.Client) AuthSelector {
	if client.Jar == nil {
		jar, _ := cookiejar.New(nil)
		// 绑定session
		client.Jar = jar
	}
	return &useCtx{Client: client, CAS: &cas{}}
}

// 选择鉴权方式
type AuthSelector interface {
	WithAuth(username, password string) PlatformSelector
	WithToken(token string) PlatformSelector
}

// 选择平台
type PlatformSelector interface {
	On(platform Platform) ActionSelector
}

// 选择要执行的动作
type ActionSelector interface {
	Login() error
	LoginService(url string) (string, error)
}

type useCtx struct {
	// 请求客户端
	Client *http.Client

	CAS *cas
}

var _ AuthSelector = &useCtx{}
var _ PlatformSelector = &useCtx{}
var _ ActionSelector = &useCtx{}

// 选择一网通平台或 Webvpn平台
func (c *useCtx) On(platform Platform) ActionSelector {
	if platform == WebVPN {
		c.CAS.Domain = webvpnDomain
		c.CAS.BaseURL = webvpnBaseURL
	} else {
		c.CAS.Domain = casDomain
		c.CAS.BaseURL = casBaseURL
	}
	return c
}

// 使用账号密码
func (c *useCtx) WithAuth(username, password string) PlatformSelector {
	c.CAS.UseToken = false
	c.CAS.Username = username
	c.CAS.Password = password
	return c
}

// 使用Token
func (c *useCtx) WithToken(token string) PlatformSelector {
	c.CAS.UseToken = true
	c.CAS.Token = token
	return c
}

// 登陆
func (c *useCtx) Login() error {
	_, err := c.LoginService(portalURL)
	return err
}

// 登陆指定服务，url需要是服务的完整地址，例如
// https://219.216.96.4/eams/homeExt.action
// 返回页面内容，如果登陆失败会返回error
func (c *useCtx) LoginService(url string) (string, error) {
	c.CAS.ServiceURL = url
	return c.CAS.Login(c.Client)
}

/*
	About
*/

// TODO 查询
// 提供查询相关信息的链式调用
func About(client *http.Client) QuerySelector {
	return &aboutCtx{Client: client}
}

// 选择要查询的内容
type QuerySelector interface {
	Token(platform Platform) (string, error)
	Info(platform Platform) (*PersonalInfo, error)
}

type aboutCtx struct {
	Client *http.Client
}

var _ QuerySelector = &aboutCtx{}

func (c *aboutCtx) Token(platform Platform) (string, error) {
	var domain string
	if platform == WebVPN {
		domain = webvpnDomain
	} else {
		domain = casDomain
	}
	return getToken(c.Client, domain)
}

func (c *aboutCtx) Info(platform Platform) (*PersonalInfo, error) {
	return nil, nil
}