package main

import (
	"testing"
)

type CountTableTest struct {
	Name   string
	Input  []string
	Output map[string]int
}

func TestCount(t *testing.T) {
	tests := []CountTableTest{
		{
			Name:   "empty",
			Input:  []string{},
			Output: map[string]int{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res := Count(test.Input)
			if len(res) != len(test.Output) {
				t.Fail()
			}
		})
	}
}

type MySquareTest struct {
	Name   string
	Input  int
	Output int
}

func TestMySquare(t *testing.T) {
	tests := []MySquareTest{
		{
			Name:   "0",
			Input:  0,
			Output: 0,
		},
		{
			Name:   "3",
			Input:  3,
			Output: 9,
		},
		{
			Name:   "-1",
			Input:  -1,
			Output: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res := MySquare(test.Input)
			if test.Output != res {
				t.Fail()
			}
		})
	}
}
