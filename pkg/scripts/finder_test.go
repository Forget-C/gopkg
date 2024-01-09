package scripts

import (
	"context"
	"github.com/Forget-C/gopkg/pkg/log"
	"os"
	"testing"
)

func TestScan(t *testing.T) {
	type args struct {
		dir   string
		skipE bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				dir:   "/Users/chenyang/Projects/oschina/gitee-cloud-ide-deploy/cmd/pre_init/test_scripts/pre_scripts",
				skipE: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scan(tt.args.dir, tt.args.skipE)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			t.Log(got.GetExecuteScripts())
			t.Log(got.GetRollbackScripts(-1))
		})
	}
}

func TestManager_Execute(t *testing.T) {
	manager := NewScriptManager("/Users/chenyang/Projects/oschina/gitee-cloud-ide-deploy/cmd/pre_init/test_scripts", log.Logger, os.Stdout)
	err := manager.AddPreScripts("pre_scripts")
	if err != nil {
		t.Fatal(err)
	}
	err = manager.RunPrepare()
	if err != nil {
		t.Fatal(err)
	}
	manager.Execute(context.Background())
}

func Test_getNameNumber(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "001_init_rollback.sh",
			args: args{
				name: "001_init_rollback.sh",
			},
			want: 1,
		},
		{
			name: "0100_init.sh",
			args: args{
				name: "001_init.sh",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNameNumber(tt.args.name); got != tt.want {
				t.Errorf("getNameNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRollback(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "001_init_rollback.sh",
			args: args{
				name: "001_init_rollback.sh",
			},
			want: true,
		},
		{
			name: "0100_init.sh",
			args: args{
				name: "001_init.sh",
			},
			want: false,
		},
		{
			name: "0100_init-rollback.sh",
			args: args{
				name: "001_init-rollback.sh",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRollback(tt.args.name); got != tt.want {
				t.Errorf("isRollback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "001_init_rollback.sh",
			args: args{
				name: "001_init_rollback.sh",
			},
			want: "init",
		},
		{
			name: "0100_init.sh",
			args: args{
				name: "001_init.sh",
			},
			want: "init",
		},
		{
			name: "0100_init-rollback.sh",
			args: args{
				name: "001_init-rollback.sh",
			},
			want: "init-rollback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getName(tt.args.name); got != tt.want {
				t.Errorf("getName() = %v, want %v", got, tt.want)
			}
		})
	}
}
