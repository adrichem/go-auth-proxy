package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Proxy(Upstream string, HeaderName string, HeaderValue string) gin.HandlerFunc {
	remote, err := url.Parse(Upstream)
	if nil != err {
		panic(err)
	}

	if remote.Scheme == "" {
		panic("missing scheme on upstream url")
	}
	if remote.Host == "" {
		panic("missing hosy on upstream url")
	}

	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			if HeaderName != "" {
				req.Header.Add(HeaderName, HeaderValue)
			}
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("path")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
