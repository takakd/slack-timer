package config

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cases := []struct {
		name         string
		key          string
		defaultValue string
		value        string
		want         string
	}{
		{name: "ok:no params", key: "test1", defaultValue: "other1", value: "value1", want: "value1"},
		{name: "ok:params", key: "test2", defaultValue: "value2", value: "value2", want: "value2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := NewMockConfig(ctrl)
			m.EXPECT().Get(gomock.Eq(c.key), gomock.Eq(c.defaultValue)).Return(c.value)

			SetConfig(m)

			assert.Equal(t, c.want, Get(c.key, c.defaultValue))
		})
	}
}
