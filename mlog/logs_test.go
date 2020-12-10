package mlog

import "testing"

func TestLog(t *testing.T) {

	SetLevel(DEBUG)
	Errorf("",ERROR)
}
