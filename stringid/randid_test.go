package stringid_test

import (
	"fmt"
	"github.com/sisukas/commons/stringid"
	"testing"
)

func TestCreatingID(t *testing.T) {
	//const len = 1000000
	const len = 100
	idmap := make(map[string]int)

	for i := 0; i < len; i++ {
		id := stringid.RandID(8)
		_, exists := idmap[id]
		if exists {
			t.Errorf("ID collission %s", id)
			continue
		}
		idmap[id] = i
		fmt.Printf("\n%s", id)
	}
}

func TestProfanity(t *testing.T) {

	if !stringid.IsProfanity("fuck") {
		t.Error("Profanity not detected ")
	}
}
