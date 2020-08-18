package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"proteinreminder/internal/httputil"
	"proteinreminder/internal/ioc"
	"proteinreminder/internal/log"
	"proteinreminder/internal/controller"
)

const (
	ApiPrefixPath     = "/api"
	Version           = "1.0"
	DefaultServerPort = "8080"
)

// --------------------------------------------------------

// Controllers implement request handlers according to this type.
type WithContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

type handlerFunc func(w http.ResponseWriter, r *http.Request)

// Handler with access log.
// Ref: https://golang.org/doc/articles/wiki/
// Ref: https://ema-hiro.hatenablog.com/entry/2018/05/14/003526
func makeHandlerFunc(ctx context.Context, logger log.Logger, f WithContextHandlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Output access log.
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		logger.Info(fmt.Sprintf("Remote: %s [%s] %s\n", rAddr, method, path))

		// Call actions.
		f(ctx, w, r)
	}
}

// --------------------------------------------------------

type Server struct {
	addr string
}

func NewServer() *Server {
	s := &Server{}
	if port := os.Getenv("PORT"); port != "" {
		s.addr = ":" + port
	} else {
		s.addr = ":" + DefaultServerPort
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	logger := ioc.GetLogger()

	// POST: /api/<ver>/slack-callback
	http.HandleFunc("/"+Version+"/test", makeHandlerFunc(ctx, logger, controller.SlackCallbackHandler))

	// GET: /api/<ver>/test
	http.HandleFunc(ApiPrefixPath+"/"+Version+"/test", makeHandlerFunc(ctx, logger, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		httputil.WriteJsonResponse(w, 200, []byte(fmt.Sprintf("called /%s/test", Version)))
	}))

	return http.ListenAndServe(s.addr, nil)
}
