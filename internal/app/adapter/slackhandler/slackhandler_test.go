package slackhandler

import (
	"errors"
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/util/di"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSlackHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := slack.NewMockAPI(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get(gomock.Eq("slack.API")).Return(s)
	di.SetDi(d)

	h := NewSlackHandler()
	assert.Equal(t, s, h.api)
}

func TestSlackApi_Notify(t *testing.T) {
	caseUserID := "test user"
	caseMessage := "test message"
	caseChannelID := "test channel id"

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockAPI(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserID)).Return(caseChannelID, nil)
		s.EXPECT().ChatPostMessage(gomock.Eq(caseChannelID), gomock.Eq(caseMessage)).Return(nil)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slack.API")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserID, caseMessage)
		assert.NoError(t, err)
	})

	t.Run("ng:conversations.open", func(t *testing.T) {
		caseError := errors.New("test error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockAPI(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserID)).Return("", caseError)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slack.API")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserID, caseMessage)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:chatpostmessage", func(t *testing.T) {
		caseError := errors.New("test error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockAPI(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserID)).Return(caseChannelID, nil)
		s.EXPECT().ChatPostMessage(gomock.Eq(caseChannelID), gomock.Eq(caseMessage)).Return(caseError)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slack.API")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserID, caseMessage)
		assert.Equal(t, caseError, err)
	})
}
