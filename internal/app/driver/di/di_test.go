package di

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"proteinreminder/internal/pkg/config"
	"testing"
)

type UnitTestDi struct {
}

func (t *UnitTestDi) Get(name string) interface{} {
	if name == "test" {
		return "value"
	}
	return nil
}

func TestGet(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		SetDi(&UnitTestDi{})
		v := di.Get("test")
		assert.Equal(t, "value", v.(string))
	})
}

func TestSetConfig(t *testing.T) {
	cases := []struct {
		name  string
		env   string
		value interface{}
	}{
		{"production", "production", nil},
		{"test", "test", &TestDi{}},
		{"empty", "", nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := config.NewMockConfig(ctrl)
			m.EXPECT().Get("APP_ENV", "dev").Return(c.env)
			config.SetConfig(m)
			setDi()
			assert.Equal(t, c.value, di)
		})
	}
}
