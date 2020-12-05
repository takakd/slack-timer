package settime

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/adapter/slackcontroller"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestLambdaInput_HandlerInput(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		data := slackcontroller.EventCallbackData{
			Token:  "test",
			TeamId: "test id",
			MessageEvent: slackcontroller.MessageEvent{
				Type:    "message",
				EventTs: "1234.0000001",
				User:    "YIG35ADg",
				Ts:      "1234.0000001",
				Text:    "message",
			},
			Challenge: "challenge",
		}
		dataJson, _ := json.Marshal(&data)

		caseInput := LambdaInput{
			Body: string(dataJson),
		}

		want := &slackcontroller.HandlerInput{
			EventData: data,
		}

		got, err := caseInput.HandlerInput()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ng", func(t *testing.T) {
		caseInput := LambdaInput{
			Body: "{syntax error",
		}

		got, err := caseInput.HandlerInput()
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestLambdaHandler(t *testing.T) {
	t.Run("ok:struct response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		data := slackcontroller.HandlerInput{
			EventData: slackcontroller.EventCallbackData{
				Token:  "test",
				TeamId: "test id",
				MessageEvent: slackcontroller.MessageEvent{
					Type:    "message",
					EventTs: "1234.0000001",
					User:    "YIG35ADg",
					Ts:      "1234.0000001",
					Text:    "message",
				},
				Challenge: "challenge",
			},
		}
		dataJson, _ := json.Marshal(&data.EventData)

		caseInput := LambdaInput{
			Body: string(dataJson),
		}

		caseRespBody := struct {
			Message string `json:"message"`
		}{
			"test",
		}
		caseResp := slackcontroller.Response{
			StatusCode: http.StatusOK,
			Body:       caseRespBody,
		}

		mh := slackcontroller.NewMockHandler(ctrl)
		mh.EXPECT().Handler(gomock.Eq(ctx), gomock.Eq(data)).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.Handler")).Return(mh)
		di.SetDi(md)

		h := NewSetTimerLambdaHandler()
		got, err := h.LambdaHandler(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, true, got.IsBase64Encoded)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		wantRespBody, _ := json.Marshal(caseRespBody)
		wantResp := &LambdaOutput{
			IsBase64Encoded: true,
			StatusCode:      http.StatusOK,
			Body:            string(wantRespBody),
		}
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok:string response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		data := slackcontroller.HandlerInput{
			EventData: slackcontroller.EventCallbackData{
				Token:  "test",
				TeamId: "test id",
				MessageEvent: slackcontroller.MessageEvent{
					Type:    "message",
					EventTs: "1234.0000001",
					User:    "YIG35ADg",
					Ts:      "1234.0000001",
					Text:    "message",
				},
				Challenge: "challenge",
			},
		}
		dataJson, _ := json.Marshal(&data.EventData)

		caseInput := LambdaInput{
			Body: string(dataJson),
		}

		caseRespBody := "message"
		caseResp := slackcontroller.Response{
			StatusCode: http.StatusOK,
			Body:       caseRespBody,
		}

		mh := slackcontroller.NewMockHandler(ctrl)
		mh.EXPECT().Handler(gomock.Eq(ctx), gomock.Eq(data)).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.Handler")).Return(mh)
		di.SetDi(md)

		h := NewSetTimerLambdaHandler()
		got, err := h.LambdaHandler(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, true, got.IsBase64Encoded)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		wantResp := &LambdaOutput{
			IsBase64Encoded: true,
			StatusCode:      http.StatusOK,
			Body:            caseRespBody,
		}
		assert.Equal(t, wantResp, got)
	})

	t.Run("ng:input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		caseInput := LambdaInput{
			Body: "{invalid format",
		}

		mh := slackcontroller.NewMockHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.Handler")).Return(mh)
		di.SetDi(md)

		h := NewSetTimerLambdaHandler()
		got, err := h.LambdaHandler(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, true, got.IsBase64Encoded)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		wantResp := &LambdaOutput{
			IsBase64Encoded: true,
			StatusCode:      http.StatusInternalServerError,
			Body:            `{"message":"invalid request", "detail":"parameters are wrong"}`,
		}
		assert.Equal(t, wantResp, got)
	})
}
