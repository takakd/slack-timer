package collection_test

import (
	"fmt"
	"slacktimer/internal/pkg/collection"
)

func ExampleSet() {
	s := collection.NewSet()

	// Set values
	s.Set("value1")
	s.Set("value2")

	// Check
	fmt.Println(s.Contains("value1"))
	fmt.Println(s.Contains("value2"))

	// Remove
	s.Remove("value2")

	// Removed so does not exists.
	fmt.Println(s.Contains("value2"))

	// Output:
	// true
	// true
	// false
}
