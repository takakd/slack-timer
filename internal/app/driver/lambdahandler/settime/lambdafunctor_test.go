package settime

import (
	"context"
	"encoding/json"
	"net/http"
	"slacktimer/internal/app/adapter/settime"
	"slacktimer/internal/app/util/di"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLambdaHandler(t *testing.T) {
	t.Run("ok:struct response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		data := settime.HandleInput{
			EventData: settime.EventCallbackData{
				Token:  "test",
				TeamID: "test id",
				MessageEvent: settime.MessageEvent{
					Type:    "message",
					EventTs: "1234.0000001",
					User:    "YIG35ADg",
					Ts:      "1234.0000001",
					Text:    "message",
				},
				Challenge: "challenge",
			},
		}
		dataJSON, _ := json.Marshal(&data.EventData)

		caseInput := LambdaInput{
			Body: string(dataJSON),
		}

		caseRespBody := struct {
			Message string `json:"message"`
		}{
			"test",
		}
		caseResp := settime.Response{
			StatusCode: http.StatusOK,
			Body:       caseRespBody,
		}

		mh := settime.NewMockControllerHandler(ctrl)
		mh.EXPECT().Handle(gomock.Eq(ctx), gomock.Eq(data)).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.ControllerHandler")).Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
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

		data := settime.HandleInput{
			EventData: settime.EventCallbackData{
				Token:  "test",
				TeamID: "test id",
				MessageEvent: settime.MessageEvent{
					Type:    "message",
					EventTs: "1234.0000001",
					User:    "YIG35ADg",
					Ts:      "1234.0000001",
					Text:    "message",
				},
				Challenge: "challenge",
			},
		}
		dataJSON, _ := json.Marshal(&data.EventData)

		caseInput := LambdaInput{
			Body: string(dataJSON),
		}

		caseRespBody := "message"
		caseResp := settime.Response{
			StatusCode: http.StatusOK,
			Body:       caseRespBody,
		}

		mh := settime.NewMockControllerHandler(ctrl)
		mh.EXPECT().Handle(gomock.Eq(ctx), gomock.Eq(data)).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.ControllerHandler")).Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
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

		mh := settime.NewMockControllerHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.ControllerHandler")).Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
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
