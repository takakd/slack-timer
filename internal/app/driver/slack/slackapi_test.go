package slack

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"slacktimer/internal/pkg/config"
	"testing"
)

func TestNewSlackApiDriver(t *testing.T) {
	want := &SlackApiDriver{}
	got := NewSlackApiDriver()
	assert.Equal(t, want, got)
}

func TestSlackApiDriver_ConversationsOpen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		wantChannelId := "test channel"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponse{
				Ok:      true,
				Channel: []string{wantChannelId},
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("dummy")
		assert.Equal(t, wantChannelId, got)
		assert.NoError(t, err)
	})

	t.Run("ng:token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewSlackApiDriver()
			d.ConversationsOpen("dummy")
		})
	})

	t.Run("ng:url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return("")
		config.SetConfig(c)

		assert.Panics(t, func() {
			d := NewSlackApiDriver()
			d.ConversationsOpen("dummy")
		})
	})

	t.Run("ng:status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponse{
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
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL + "/wrong")
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("dummy")
		assert.Empty(t, got)
		assert.Error(t, err)
	})

	// TODO: fix httputils to use interface, mockable.

	t.Run("ng:response NG", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponse{
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
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("dummy")
		assert.Empty(t, got)
		assert.Error(t, err)
	})

	t.Run("ng:response wrong channel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody := &ConversationsOpenResponse{
				Ok:      true,
				Channel: []string{"unexpected1", "unexpected2"},
			}
			resp, err := json.Marshal(respBody)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(resp))
		}))
		defer server.Close()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_BOT_TOKEN"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq(
			"SLACK_API_URL_CONVERSATIONSOPEN"), gomock.Eq("")).Return(server.URL)
		config.SetConfig(c)

		d := NewSlackApiDriver()
		got, err := d.ConversationsOpen("dummy")
		assert.Empty(t, got)
		assert.Error(t, err)
	})
}

//func TestSlackApiDriver_ChatPostMessage(t *testing.T) {
//
//}
//
//func TestPostJson(t *testing.T) {
//
//}
