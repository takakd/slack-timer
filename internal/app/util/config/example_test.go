package config_test

import (
	"fmt"
	"os"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/config"
)

func ExampleGet() {
	os.Setenv("EXAMPLE_KEY", "EXAMPLE_VALUE")
	appinitializer.AppInit()

	fmt.Println(config.Get("EXAMPLE_KEY", "default"))

	// Output:
	// EXAMPLE_VALUE
}
