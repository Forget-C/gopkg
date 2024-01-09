package scripts

import (
	"context"
	"path/filepath"
	"testing"
)

func TestExecutorImpl_Run(t *testing.T) {
	dir := "/Users/chenyang/Projects/oschina/gitee-cloud-ide-deploy/cmd/pre_init/test_scripts/pre_scripts"
	t.Run("bash", func(t *testing.T) {
		f := filepath.Join(dir, "10_a.sh")
		e := MakeExecutor(f)
		ok, out, err := e.Run(context.Background(), nil)
		t.Log(ok, string(out), err)
		t.Log(err)
	})
	t.Run("python", func(t *testing.T) {
		f := filepath.Join(dir, "66_a.py")
		e := MakeExecutor(f)
		ok, out, err := e.Run(context.Background(), nil)
		t.Log(ok, string(out), err)
		t.Log(err)
	})
	t.Run("binary", func(t *testing.T) {
		f := filepath.Join(dir, "99_b")
		e := MakeExecutor(f)
		ok, out, err := e.Run(context.Background(), nil)
		t.Log(ok, string(out), err)
		t.Log(err)
	})
}
