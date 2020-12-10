package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"slacktimer/internal/app/util/config"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIDriver(t *testing.T) {
	want := &APIDriver{}
	got := NewAPIDriver()
	assert.Equal(t, want, got)
}

func TestAPIDriver_ConversationsOpen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseUserID := "test user"
		caseChannelID := "test channel"
		caseToken := "test token"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)

			var reqBody ConversationsOpenRequestBody
			err = json.Unmarshal(buf, &reqBody)
			require.NoError(t, err)
			assert.Equal(t, reqBody.Users, caseUserID)

			respBody := &ConversationsOpenResponseBody{
				Ok:      true,
				Channel: ConversationsOpenResponseBodyChannel{caseChannelID},
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return(caseToken)
		mc.EXPECT().Get("SLACK_API_URL_CONVERSATIONSOPEN", "").Return(server.URL)
		config.SetConfig(mc)

		d := NewAPIDriver()
		got, err := d.ConversationsOpen(appcontext.TODO(), caseUserID)
		assert.Equal(t, caseChannelID, got)
		assert.NoError(t, err)
	})

	t.Run("ng:url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_URL_CONVERSATIONSOPEN", "").Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewAPIDriver()
			d.ConversationsOpen(appcontext.TODO(), "test")
		})
	})

	t.Run("ng:status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponseBody{
				Ok:    false,
				Error: "invalid token",
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		c.EXPECT().Get("SLACK_API_URL_CONVERSATIONSOPEN", "").Return(server.URL + "/wrong")
		config.SetConfig(c)

		d := NewAPIDriver()
		got, err := d.ConversationsOpen(appcontext.TODO(), "test")
		assert.Empty(t, got)
		assert.Error(t, err)
	})

	t.Run("ng:response NG", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponseBody{
				Ok:    false,
				Error: "invalid token",
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		c.EXPECT().Get("SLACK_API_URL_CONVERSATIONSOPEN", "").Return(server.URL)
		config.SetConfig(c)

		d := NewAPIDriver()
		got, err := d.ConversationsOpen(appcontext.TODO(), "test")
		assert.Empty(t, got)
		assert.Error(t, err)
	})
}

func TestAPIDriver_ChatPostMessage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseToken := "test token"
		caseChannelID := "test channel"
		caseText := "test message"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)

			var reqBody ChatPostMessageRequestBody
			err = json.Unmarshal(buf, &reqBody)
			require.NoError(t, err)
			assert.Equal(t, reqBody.Channel, caseChannelID)
			assert.Equal(t, reqBody.Text, caseText)

			respBody := &ChatPostMessageResponseBody{
				Ok: true,
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return(caseToken)
		c.EXPECT().Get("SLACK_API_URL_CHATPOSTMESSAGE", "").Return(server.URL)
		config.SetConfig(c)

		d := NewAPIDriver()
		err := d.ChatPostMessage(appcontext.TODO(), caseChannelID, caseText)
		assert.NoError(t, err)
	})

	t.Run("ng:url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_URL_CHATPOSTMESSAGE", "").Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewAPIDriver()
			d.ChatPostMessage(appcontext.TODO(), "test", "test")
		})
	})

	t.Run("ng:status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ChatPostMessageResponseBody{
				Ok: false,
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		c.EXPECT().Get("SLACK_API_URL_CHATPOSTMESSAGE", "").Return(server.URL + "/wrong")
		config.SetConfig(c)

		d := NewAPIDriver()
		err := d.ChatPostMessage(appcontext.TODO(), "test", "test")
		assert.Error(t, err)
	})

	t.Run("ng:response NG", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponseBody{
				Ok: false,
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		c.EXPECT().Get("SLACK_API_URL_CHATPOSTMESSAGE", "").Return(server.URL)
		config.SetConfig(c)

		d := NewAPIDriver()
		err := d.ChatPostMessage(appcontext.TODO(), "test", "test")
		assert.Error(t, err)
	})
}

func TestPostJson(t *testing.T) {
	t.Run("ng:marshal", func(t *testing.T) {
		invalidBody := make(chan int)
		resp, err := postJSON("http://localhost", invalidBody)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("ng:invalid url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		config.SetConfig(c)

		// Schema error
		body := &ConversationsOpenRequestBody{}
		resp, err := postJSON("not support protocol schema url", body)
		t.Log(err)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseResponse := "{}"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Check only if the request is success.
			fmt.Fprint(w, caseResponse)
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("SLACK_API_BOT_TOKEN", "").Return("test")
		config.SetConfig(c)

		// Schema error
		body := &ConversationsOpenRequestBody{}
		resp, err := postJSON(server.URL, body)
		assert.NoError(t, err)

		got, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, string(got), caseResponse)
	})
}
