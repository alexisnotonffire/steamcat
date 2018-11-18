package steamcat

import (
	"os"
	"testing"
)

var testLib Library
var gameA, gameB, gameC Game

func TestTaggedWith(t *testing.T) {
	testCases := []struct {
		lib Library
		tag string
		exp Library
		err error
	}{
		{lib: testLib, tag: "Arcade", exp: Library{&gameA}, err: nil},
	}

	for _, tc := range testCases {
		taggedLib := tc.lib.TaggedWith(tc.tag)
		if len(taggedLib) != len(tc.exp) {
			t.Errorf("Found a different number of games than expected.\nGot: %v,\nExp: %v", taggedLib, tc.exp)
		}

	VALUE:
		for _, g := range tc.exp {
			for _, h := range taggedLib {
				if g == h {
					continue VALUE
				}
			}
			t.Errorf("Could not find expected game %s\nGot: %v\nExp: %v", g.Name, taggedLib, tc.exp)
		}
	}
}

func TestMain(m *testing.M) {

	testLib = Library{&gameA, &gameB, &gameC}
	gameA = Game{ID: 1, Name: "A", Tags: []string{"Arcade", "Sports"}}
	gameB = Game{ID: 2, Name: "B", Tags: []string{"Action", "Sports"}}
	gameC = Game{ID: 3, Name: "C", Tags: []string{}}

	retCode := m.Run()

	testLib = nil
	gameA = Game{}
	gameB = Game{}
	gameC = Game{}

	os.Exit(retCode)
}
