package slackhandler

import (
	"errors"
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSlackHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := slack.NewMockAPI(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("slack.API").Return(s)
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

		ms := slack.NewMockAPI(ctrl)
		ms.EXPECT().ConversationsOpen(appcontext.TODO(), caseUserID).Return(caseChannelID, nil)
		ms.EXPECT().ChatPostMessage(appcontext.TODO(), caseChannelID, caseMessage).Return(nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("slack.API").Return(ms)
		di.SetDi(md)

		h := NewSlackHandler()
		err := h.Notify(appcontext.TODO(), caseUserID, caseMessage)
		assert.NoError(t, err)
	})

	t.Run("ng:conversations.open", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseError := errors.New("test error")

		ms := slack.NewMockAPI(ctrl)
		ms.EXPECT().ConversationsOpen(appcontext.TODO(), caseUserID).Return("", caseError)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("slack.API").Return(ms)
		di.SetDi(md)

		h := NewSlackHandler()
		err := h.Notify(appcontext.TODO(), caseUserID, caseMessage)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:chatpostmessage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseError := errors.New("test error")

		ms := slack.NewMockAPI(ctrl)
		ms.EXPECT().ConversationsOpen(appcontext.TODO(), caseUserID).Return(caseChannelID, nil)
		ms.EXPECT().ChatPostMessage(appcontext.TODO(), caseChannelID, caseMessage).Return(caseError)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("slack.API").Return(ms)
		di.SetDi(md)

		h := NewSlackHandler()
		err := h.Notify(appcontext.TODO(), caseUserID, caseMessage)
		assert.Equal(t, caseError, err)
	})
}
