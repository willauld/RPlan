package main

import (
	"fmt"
	"testing"
)

func TestDoItWithUnMarshal(t *testing.T) {
	doItWithUnMarshal()
}

func TestGetTomlData(t *testing.T) {
	tests := []struct {
		sip int
	}{
		{
			sip: 9,
		},
	}
	for i, elem := range tests {
		if elem.sip == 11 {
			err := fmt.Errorf("my mistake")
			t.Errorf("TestGetTomlData case %d: %s", i, err)
			continue
		}
		goGetTomlData()
		try2()
	}
}
