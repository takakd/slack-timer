package settime

import (
	"errors"
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"slacktimer/internal/app/usecase/timeronevent"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOnEventOutputReceivePresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewOnEventOutputReceivePresenter()
	})
}

func TestOnEventOutputReceivePresenter_Output(t *testing.T) {
	cases := []struct {
		name string
		data *timeronevent.OutputData
	}{
		{name: "ok", data: &timeronevent.OutputData{
			Result: nil,
		}},
		{name: "error", data: &timeronevent.OutputData{
			Result: errors.New("error case"),
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ac := appcontext.TODO()

			p := NewOnEventOutputReceivePresenter()

			if c.data.Result != nil {
				ml := log.NewMockLogger(ctrl)
				ml.EXPECT().ErrorWithContext(ac, "settime offevent outputport", c.data.Result)
				log.SetDefaultLogger(ml)

				p.Output(appcontext.TODO(), *c.data)

				want := &Response{
					StatusCode: http.StatusInternalServerError,
					Body:       "internal server error",
					Error:      c.data.Result,
				}
				assert.Equal(t, want, &p.Resp)

			} else {
				var err error
				c.data.SavedEvent, err = enterpriserule.NewTimerEvent("test")
				c.data.SavedEvent.Text = "Hi!"
				require.NoError(t, err)

				p.Output(appcontext.TODO(), *c.data)

				want := &Response{
					StatusCode: http.StatusOK,
					Body:       "success",
					Error:      c.data.Result,
				}
				assert.Equal(t, want, &p.Resp)
			}

		})
	}
}
