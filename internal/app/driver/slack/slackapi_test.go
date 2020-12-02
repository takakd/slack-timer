package slack

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"slacktimer/internal/app/util/config"
	"testing"
)

func TestNewSlackApiDriver(t *testing.T) {
	want := &SlackApiDriver{}
	got := NewSlackApiDriver()
	assert.Equal(t, want, got)
}

func TestSlackApiDriver_ConversationsOpen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseUserId := "test user"
		caseChannelId := "test channel"
		caseToken := "test token"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)

			var reqBody ConversationsOpenRequestBody
			err = json.Unmarshal(buf, &reqBody)
			require.NoError(t, err)
			assert.Equal(t, reqBody.Users, caseUserId)

			respBody := &ConversationsOpenResponseBody{
				Ok:      true,
				Channel: ConversationsOpenResponseBodyChannel{caseChannelId},
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return(caseToken)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen(caseUserId)
		assert.Equal(t, caseChannelId, got)
		assert.NoError(t, err)
	})

	t.Run("ng:url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewSlackApiDriver()
			d.ConversationsOpen("test")
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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL + "/wrong")
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("test")
		assert.Empty(t, got)
		assert.Error(t, err)
	})

	// TODO: fix httputils to use interface, mockable.

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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("test")
		assert.Empty(t, got)
		assert.Error(t, err)
	})
}

func TestSlackApiDriver_ChatPostMessage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseToken := "test token"
		caseChannelId := "test channel"
		caseText := "test message"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)

			var reqBody ChatPostMessageRequestBody
			err = json.Unmarshal(buf, &reqBody)
			require.NoError(t, err)
			assert.Equal(t, reqBody.Channel, caseChannelId)
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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return(caseToken)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CHATPOSTMESSAGE"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		err := d.ChatPostMessage(caseChannelId, caseText)
		assert.NoError(t, err)
	})

	t.Run("ng:url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CHATPOSTMESSAGE"), gomock.Eq("")).Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewSlackApiDriver()
			d.ChatPostMessage("test", "test")
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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CHATPOSTMESSAGE"), gomock.Eq("")).Return(server.URL + "/wrong")
		config.SetConfig(c)

		d := NewSlackApiDriver()
		err := d.ChatPostMessage("test", "test")
		assert.Error(t, err)
	})

	// TODO: fix httputils to use interface, mockable.

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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CHATPOSTMESSAGE"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		err := d.ChatPostMessage("test", "test")
		assert.Error(t, err)
	})
}

func TestPostJson(t *testing.T) {
	t.Run("ng:marshal", func(t *testing.T) {
		invalidBody := make(chan int)
		resp, err := postJson("http://localhost", invalidBody)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("ng:invalid url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		config.SetConfig(c)

		// Schema error
		body := &ConversationsOpenRequestBody{}
		resp, err := postJson("not support protocol schema url", body)
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
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("test")
		config.SetConfig(c)

		// Schema error
		body := &ConversationsOpenRequestBody{}
		resp, err := postJson(server.URL, body)
		assert.NoError(t, err)

		got, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, string(got), caseResponse)
	})
}
