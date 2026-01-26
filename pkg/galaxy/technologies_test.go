package galaxy

import (
	"reflect"
	"testing"
)

func TestTechnologiesResearch(t *testing.T) {
	tech := NewTechnologies()

	tech.Research(100, 50, 200, 150)

	expectedTechnologies := &Technologies{
		Attack:  2,
		Defense: 1.5,
		Engine:  3,
		Cargo:   2.5,
	}

	if !reflect.DeepEqual(expectedTechnologies, tech) {
		t.Errorf("technologies are not equal %v != %v", expectedTechnologies, tech)
	}
}
