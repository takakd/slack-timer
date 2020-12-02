package slackhandler

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestNewSlackHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := slack.NewMockSlackApi(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get(gomock.Eq("slackhandler.SlackApi")).Return(s)
	di.SetDi(d)

	h := NewSlackHandler()
	assert.Equal(t, s, h.api)
}

func TestSlackApi_Notify(t *testing.T) {
	caseUserId := "test user"
	caseMessage := "test message"
	caseChannelId := "test channel id"

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockSlackApi(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserId)).Return(caseChannelId, nil)
		s.EXPECT().ChatPostMessage(gomock.Eq(caseChannelId), gomock.Eq(caseMessage)).Return(nil)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slackhandler.SlackApi")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserId, caseMessage)
		assert.NoError(t, err)
	})

	t.Run("ng:conversations.open", func(t *testing.T) {
		caseError := errors.New("test error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockSlackApi(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserId)).Return("", caseError)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slackhandler.SlackApi")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserId, caseMessage)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:chatpostmessage", func(t *testing.T) {
		caseError := errors.New("test error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := slack.NewMockSlackApi(ctrl)
		s.EXPECT().ConversationsOpen(gomock.Eq(caseUserId)).Return(caseChannelId, nil)
		s.EXPECT().ChatPostMessage(gomock.Eq(caseChannelId), gomock.Eq(caseMessage)).Return(caseError)

		d := di.NewMockDI(ctrl)
		d.EXPECT().Get(gomock.Eq("slackhandler.SlackApi")).Return(s)
		di.SetDi(d)

		h := NewSlackHandler()
		err := h.Notify(caseUserId, caseMessage)
		assert.Equal(t, caseError, err)
	})
}
