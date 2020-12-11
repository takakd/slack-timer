package log_test

import (
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/log"
)

func ExampleInfo() {
	appinitializer.AppInit()

	v := 123
	log.Info("info log", v)

	f := 123.456
	log.Debug("debug log", f)

	s := "some string"
	log.Error("error log", s)

	// Output:
	// {"level":"INFO","msg":["info log",123]}
	// {"level":"DEBUG","msg":["debug log",123.456]}
	// {"level":"ERROR","msg":["error log","some string"]}
}

func ExampleErrorWithContext() {
	appinitializer.AppInit()
	ac := appcontext.TODO()

	v := 123
	log.ErrorWithContext(ac, "info log", v)

	f := 123.456
	log.ErrorWithContext(ac, "info log", f)

	s := "some string"
	log.ErrorWithContext(ac, "info log", s)

	// Output:
	// {"AwsRequestID":"","level":"ERROR","msg":["info log",123]}
	// {"AwsRequestID":"","level":"ERROR","msg":["info log",123.456]}
	// {"AwsRequestID":"","level":"ERROR","msg":["info log","some string"]}
}
