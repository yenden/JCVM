package test

import (
	"core"
	"testing"
)

func TestBuildApplet(t *testing.T) {
	var app *core.CardApplet
	args := []string{`framework.ijc`, `lang.ijc`, `helloword.ijc`}
	for i := 0; i < len(args)-1; i++ {
		dataBuffer := core.ReadInBuffer(args[i])
		app = core.BuildApplet(dataBuffer, len(dataBuffer))
		if app == nil {
			t.Errorf("%s parsing return nil", args[i])
		}
	}

}
