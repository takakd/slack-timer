package settime

import (
	"errors"
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSaveEventOutputReceivePresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewSaveEventOutputReceivePresenter()
	})
}

func TestSaveEventOutputReceivePresenter_Output(t *testing.T) {
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

			ac := appcontext.TODO()

			p := NewSaveEventOutputReceivePresenter()

			if c.data.Result != nil {
				ml := log.NewMockLogger(ctrl)
				ml.EXPECT().ErrorWithContext(ac, "settime outputport", c.data.Result)
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
