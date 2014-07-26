package stringset

import (
	"reflect"
	"testing"
)

func TestIndexof(t *testing.T) {

	var tests = []struct {
		list []string
		item string
		want int
	}{
		{[]string{"hello"}, "hello", 0},
		{[]string{"hello"}, "bello", -1},
		{[]string{"hello", "bello"}, "bello", 1},
		{[]string{}, "bello", -1},
	}
	for _, tc := range tests {
		if got := IndexOf(tc.list, tc.item); got != tc.want {
			t.Errorf("In slice %v, %v should be at position %d but got %d", tc.list, tc.item, tc.want, got)
		}
	}
}

func TestAdd(t *testing.T) {

	var tests = []struct {
		list []string
		item string
		want []string
	}{
		{[]string{"hello"}, "hello", []string{"hello"}},
		{[]string{"hello"}, "bello", []string{"hello", "bello"}},
		{[]string{"hello", "bello"}, "bello", []string{"hello", "bello"}},
		{[]string{}, "bello", []string{"bello"}},
	}
	for i, tc := range tests {
		if got := Add(tc.list, tc.item); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("case %d failed: want %v, got %v", i, tc.want, got)
		}
	}
}

func TestRemove(t *testing.T) {

	var tests = []struct {
		list []string
		item string
		want []string
	}{
		{[]string{"hello"}, "hello", []string{}},
		{[]string{"hello"}, "bello", []string{"hello"}},
		{[]string{"hello", "bello"}, "bello", []string{"hello"}},
		{[]string{}, "bello", []string{}},
	}
	for i, tc := range tests {
		if got := Remove(tc.list, tc.item); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("case %d failed: want %v, got %v", i, tc.want, got)
		}
	}
}
