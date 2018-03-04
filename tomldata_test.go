package main

import (
	"fmt"
	"testing"

	"github.com/willauld/rplanlib"
)

func TestDoItWithUnMarshal(t *testing.T) {
	doItWithUnMarshal()
}

func TestGetTomlData(t *testing.T) {
	tests := []struct {
		toml string
	}{
		{
			toml: "mobile_j.toml",
		},
		{
			toml: "hack.toml",
		},
	}
	for i, elem := range tests {
		fmt.Printf("------ Case %d -----------\n", i)
		//goGetTomlData()
		ms := getInputStringsMapFromToml(elem.toml)
		if ms == nil {
			t.Errorf("TestGetTomlData case %d: ms is nil", i)
			continue
		}
		for _, v := range rplanlib.InputStrDefs {
			r, ok := ms[v]
			if !ok {
				t.Errorf("TestGetTomlData case %d: missing ms[%s]", i, v)
			}
			if r != "" {
				fmt.Printf("    %s: '%s'\n", v, r)
			}
		}
		for x := 1; x < rplanlib.MaxStreams+1; x++ {
			for _, v := range rplanlib.InputStreamStrDefs {
				r, ok := ms[fmt.Sprintf("%s%d", v, x)]
				if !ok {
					t.Errorf("TestGetTomlData case %d: missing ms[%s]", i, v)
				}
				if r != "" {
					fmt.Printf("    %s%d: '%s'\n", v, x, r)
				}
			}
		}
	}
}
