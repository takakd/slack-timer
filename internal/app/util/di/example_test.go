package di_test

import (
	"fmt"
	"os"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
)

func ExampleGet() {
	os.Setenv("APP_ENV", "dev")
	appinitializer.AppInit()

	s := di.Get("settime.OnEventOutputReceivePresenter")
	fmt.Printf("%T", s)

	// Output:
	// *settime.OnEventOutputReceivePresenter
}
