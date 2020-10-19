package webserver

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/log"
	"strings"
	"testing"
	"time"
)

func TestMakeHandlerFunc(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		req.RemoteAddr = "111.222.333.444:1234"

		m := log.NewMockLogger(ctrl)
		log.SetDefaultLogger(m)
		m.EXPECT().Print(gomock.Eq(fmt.Sprintf("[INFO] Remote: %s [%s] /test\n", req.RemoteAddr, req.Method)))

		ctx := context.TODO()
		called := false
		resp := httptest.NewRecorder()
		f := makeHandlerFunc(ctx, func(c context.Context, w http.ResponseWriter, r *http.Request) {
			called = true
			assert.Equal(t, ctx, c)
			assert.Equal(t, req, r)
			assert.Equal(t, resp, w)
		})
		f(resp, req)
		assert.Equal(t, true, called)
	})
}

func TestNewWebServer(t *testing.T) {
	cases := []struct {
		name string
		port string
		want string
	}{
		{"default port", "", ":" + DefaultServerPort},
		{"set port", "1234", ":1234"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := config.NewMockConfig(ctrl)
			m.EXPECT().Get(gomock.Eq("PORT"), gomock.Eq("8080")).Return(c.port)
			config.SetConfig(m)

			ctx := context.TODO()

			s := NewWebServer(ctx)
			assert.NotNil(t, s)
			assert.Equal(t, s.server.Addr, c.want)
		})
	}
}

func TestWebServer_Run(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := config.NewMockConfig(ctrl)
		m.EXPECT().Get(gomock.Eq("PORT"), gomock.Eq("8080")).Return("8080")
		config.SetConfig(m)

		ctx := context.TODO()

		s := NewWebServer(ctx)
		assert.NotNil(t, s)

		run := make(chan bool)
		go func() {
			err := s.Run()
			assert.Equal(t, err, http.ErrServerClosed)
			run <- true

		}()

		// Wait for the server to start
		time.Sleep(2 * time.Second)

		s.server.Close()

		// Wait for the server to return the result of Run()
		isRun := <-run
		assert.Equal(t, isRun, true)
	})
}
