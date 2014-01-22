package task

import (
	"testing"
)

func TestVars(t *testing.T) {
	env := NewEnvironment()

	env.SetString("hello", "world")
	if env.GetString("hello") != "world" {
		t.Fail()
	}

	env.SetInt("numz", 123)
	if env.GetInt("numz") != 123 {
		t.Fail()
	}

	env.SetBool("yupyupyup", true)
	if env.GetBool("yupyupyup") != true {
		t.Fail()
	}

	if env.Exists("doesntexist") == true {
		t.Fail()
	}
}
