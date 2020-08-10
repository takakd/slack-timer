package server

import (
	"os"
	"net/http"
	"proteinreminder/internal/log"
	"proteinreminder/internal/ioc"
	"proteinreminder/internal/httputil"
	"fmt"
)

const (
	ApiPrefixPath = "/api"
	Version = "1.0"
)

// App Web Server
type Server struct {
	addr string
}

func NewServer() *Server {
	s := &Server{
		addr: ":8080",
	}
	if port := os.Getenv("PORT"); port != "" {
		s.addr = ":" + port
	}
	return s
}


// ref: https://ema-hiro.hatenablog.com/entry/2018/05/14/003526
func logHandlerFunc(logger log.Logger, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		logger.Info(fmt.Sprintf("Remote: %s [%s] %s\n", rAddr, method, path))
		f(w, r)
	}
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) Run() error {
	logger := ioc.GetLogger()
	//http.HandleFunc("/"+Version+"/test", logHandlerFunc(logger, controller.SlackCallbackHandler))

	http.HandleFunc(ApiPrefixPath + "/" + Version+"/test", logHandlerFunc(logger, func(w http.ResponseWriter, r *http.Request) {
		httputil.WriteJsonResponse(w, 200, []byte(fmt.Sprintf("called /%s/test", Version)))
	}))

	return http.ListenAndServe(s.addr, nil)
}
