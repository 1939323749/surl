package main

import (
	"surl/handler"
	"testing"
)

func TestValidURL(t *testing.T) {
	resultValid := handler.ValidUrl("https://www.google.com.hk")
	if !resultValid {
		t.Errorf("ValidUrl() was incorrect, got: %t, want: %t.", resultValid, true)
	}
	resultInvalid := handler.ValidUrl("123")
	if resultInvalid {
		t.Errorf("ValidUrl() was incorrect, got: %t, want: %t.", resultValid, true)
	}
}
