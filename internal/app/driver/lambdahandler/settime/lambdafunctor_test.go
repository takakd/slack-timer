package settime

import (
	"context"
	"encoding/json"
	"net/http"
	"slacktimer/internal/app/adapter/settime"
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type AppContextMatcher struct {
	testValue appcontext.AppContext
}

func (m *AppContextMatcher) String() string {
	return fmt.Sprintf("%v", m.testValue)
}
func (m *AppContextMatcher) Matches(x interface{}) bool {
	another, _ := x.(appcontext.AppContext)
	matched := true
	matched = matched && m.testValue.RequestID() == another.RequestID()
	return matched
}

func TestLambdaHandler(t *testing.T) {
	t.Run("ok:struct response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := &lambdacontext.LambdaContext{}
		ctx := lambdacontext.NewContext(context.TODO(), lc)
		ac, _ := appcontext.NewLambdaAppContext(ctx, time.Now())

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
		matcher := &AppContextMatcher{
			testValue: ac,
		}
		mh.EXPECT().Handle(matcher, data).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("settime.ControllerHandler").Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, false, got.IsBase64Encoded)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		wantRespBody, _ := json.Marshal(caseRespBody)
		wantResp := &LambdaOutput{
			IsBase64Encoded: false,
			StatusCode:      http.StatusOK,
			Body:            string(wantRespBody),
		}
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok:string response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := &lambdacontext.LambdaContext{}
		ctx := lambdacontext.NewContext(context.TODO(), lc)
		ac, _ := appcontext.NewLambdaAppContext(ctx, time.Now())

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

		matcher := &AppContextMatcher{
			testValue: ac,
		}
		mh.EXPECT().Handle(matcher, data).Return(&caseResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("settime.ControllerHandler").Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, false, got.IsBase64Encoded)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		wantResp := &LambdaOutput{
			IsBase64Encoded: false,
			StatusCode:      http.StatusOK,
			Body:            caseRespBody,
		}
		assert.Equal(t, wantResp, got)
	})

	t.Run("ng:input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := &lambdacontext.LambdaContext{}
		ctx := lambdacontext.NewContext(context.TODO(), lc)

		caseInput := LambdaInput{
			Body: "{invalid format",
		}

		mh := settime.NewMockControllerHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("settime.ControllerHandler").Return(mh)
		di.SetDi(md)

		h := NewLambdaFunctor()
		got, err := h.Handle(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, false, got.IsBase64Encoded)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		wantResp := &LambdaOutput{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            `{"message":"invalid request", "detail":"parameters are wrong"}`,
		}
		assert.Equal(t, wantResp, got)
	})
}
