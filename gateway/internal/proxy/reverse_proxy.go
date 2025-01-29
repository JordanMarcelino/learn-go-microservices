package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/log"
)

func NewReverseProxy(target string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Logger.Fatalf("failed parsing url: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(ctx *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				ctx.Error(err)
				ctx.Abort()
			}
		}()

		params := map[string]any{
			"path":   ctx.Request.URL.Path,
			"target": target,
		}

		log.Logger.WithFields(params).Info("proxying request")

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = ctx.Param("path")
			req.Header = ctx.Request.Header
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
