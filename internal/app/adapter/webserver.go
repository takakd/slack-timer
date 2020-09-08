package adapter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
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
	addr string
}

func NewWebServer() *WebServer {
	s := &WebServer{
		addr: ":" + DefaultServerPort,
	}
	if port := os.Getenv("PORT"); port != "" {
		s.addr = ":" + port
	}
	return s
}

// Run server process.
func (s *WebServer) Run(ctx context.Context) error {

	// POST: /api/<ver>/slack-callback
	http.HandleFunc("/"+Version+"/test", makeHandlerFunc(ctx, SlackCallbackHandler))

	// GET: /api/<ver>/test
	http.HandleFunc(ApiPrefixPath+"/"+Version+"/test", makeHandlerFunc(ctx, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		httputil.WriteJsonResponse(w, 200, []byte(fmt.Sprintf("called /%s/test", Version)))
	}))

	return http.ListenAndServe(s.addr, nil)
}
