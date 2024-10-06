/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>

*/

package cmd_test

import (
	// Use this if testing errors

	"testing"
)

// Example function to test (You need to replace this with your actual function)
func MyFunction(a, b int) (int, error) {
	// Placeholder: Replace with actual implementation
	return 0, nil
}

// Basic Test Function
func TestMyFunction_Basic(t *testing.T) {
	// Example with no errors expected
	result, err := MyFunction(2, 3)
	expected := 0

	if err != nil {
		t.Fatalf("MyFunction(2, 3) returned an unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("MyFunction(2, 3) failed. Expected %d, got %d", expected, result)
	}
}
