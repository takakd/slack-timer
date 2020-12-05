package slackcontroller

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/log"
	"testing"
)

func TestNewLambdaResponseAdaptPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewLambdaResponseAdaptPresenter()
	})
}

func TestLambdaResponseAdaptPresenter_Output(t *testing.T) {
	cases := []struct {
		name string
		data *updatetimerevent.OutputData
	}{
		{name: "ok", data: &updatetimerevent.OutputData{
			Result: nil,
		}},
		{name: "error", data: &updatetimerevent.OutputData{
			Result: errors.New("error case"),
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := NewLambdaResponseAdaptPresenter()

			if c.data.Result != nil {
				ml := log.NewMockLogger(ctrl)
				ml.EXPECT().Info(gomock.Eq(fmt.Sprintf("SetRequestOutputPort.Output error=%v", c.data.Result)))
				log.SetDefaultLogger(ml)

				p.Output(*c.data)

				want := &Response{
					StatusCode: http.StatusInternalServerError,
					Body:       "internal server error",
					Error:      c.data.Result,
				}
				assert.Equal(t, want, p.Resp)

			} else {
				var err error
				c.data.SavedEvent, err = enterpriserule.NewTimerEvent("test")
				require.NoError(t, err)

				p.Output(*c.data)

				want := &Response{
					StatusCode: http.StatusOK,
					Body:       "success",
					Error:      c.data.Result,
				}
				assert.Equal(t, want, p.Resp)
			}

		})
	}
}
