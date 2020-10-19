package webserver

import (
	"context"
	"fmt"
	"net/http"
	"proteinreminder/internal/app/adapter/slackcontroller"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
	"time"
)

const (
	ApiPrefixPath     = "/api"
	Version           = "1.0"
	DefaultServerPort = "8080"
)

// Controllers implement request handlers according to this type.
type WithContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

/*
Make server handler func with access log.
Ref: https://golang.org/doc/articles/wiki/
Ref: https://ema-hiro.hatenablog.com/entry/2018/05/14/003526
*/
func makeHandlerFunc(ctx context.Context, f WithContextHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Output access log.
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Info(fmt.Sprintf("Remote: %s [%s] %s\n", rAddr, method, path))

		// Call actions.
		f(ctx, w, r)
	}
}

/*
Web server

Initialize routing and run server process.
*/
type WebServer struct {
	server *http.Server
}

func NewWebServer(ctx context.Context) *WebServer {

	addr := ":" + DefaultServerPort
	if port := config.Get("PORT", "8080"); port != "" {
		addr = ":" + port
	}

	mux := http.NewServeMux()

	// POST: /api/<ver>/slack-callback
	mux.HandleFunc(fmt.Sprintf("%s/%s/slack-callback", ApiPrefixPath, Version), makeHandlerFunc(ctx, slackcontroller.Handler))

	// GET: /api/<ver>/test
	mux.HandleFunc(ApiPrefixPath+"/"+Version+"/test", makeHandlerFunc(ctx, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		httputil.WriteJsonResponse(w, 200, []byte(fmt.Sprintf("called /%s/test", Version)))
	}))

	s := &WebServer{
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return s
}

// Run server process.
func (s *WebServer) Run() error {
	return s.server.ListenAndServe()
}
