package entity

import (
	"proteinreminder/internal/testutil"
	"reflect"
	"testing"
)

// --------------------------------------------------------
// Repository Role

func TestGetProteinEvent(t *testing.T) {

}

func TestFindProteinEventByTime(t *testing.T) {

}

func TestSaveProteinEvent(t *testing.T) {

}

// --------------------------------------------------------
// Entity

func TestNewProteinEvent(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  *ProteinEvent
	}{
		{name: "ok", in: "id1234", out: &ProteinEvent{userId: "id1234"}},
		{name: "ng", in: "", out: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := NewProteinEvent(c.in)
			if !reflect.DeepEqual(got, c.out) {
				t.Error(testutil.MakeTestMessageWithGotWant(got, c.out))
			}
		})
	}
}
