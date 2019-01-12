package builtin

import (
	"os"
	"testing"
)

func Test_reduce(t *testing.T) {
	must(t, "/", reduce("/"))
	must(t, "/", reduce("//"))
	must(t, "/", reduce("///"))
	must(t, "/foo", reduce("/foo"))
	must(t, "/foo", reduce("//foo"))
	must(t, "/", reduce("/foo//"))
	must(t, "/boo", reduce("/foo//boo"))
	must(t, os.Getenv("HOME")+"/boo", reduce("~/boo"))
}

func must(t *testing.T, exp, got string) {
	if exp != got {
		t.Errorf("unexpected result, exp=%s got=%s", exp, got)
	}
}
