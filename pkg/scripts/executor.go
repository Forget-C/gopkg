package scripts

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Executor interface {
	// Run run script and return success, output, error
	// if out is nil, output will be return in []byte
	Run(ctx context.Context, out io.Writer) (bool, []byte, error)
	// MakeCMD make exec.Cmd
	// some executor don't need this method
	MakeCMD(ctx context.Context, out io.Writer) *exec.Cmd
	// ScriptType return script type
	ScriptType() ScriptType
}

type ScriptType string

const (
	Python ScriptType = "python"
	Bash   ScriptType = "bash"
	Binary ScriptType = "binary"
	SQL    ScriptType = "sql"
)

type ExecutorImpl struct {
	ShaBang    string
	ScriptPath string
	ScriptArgs []string
	scriptType ScriptType
}

func (e *ExecutorImpl) Run(ctx context.Context, out io.Writer) (bool, []byte, error) {
	consoleOut := true
	_out := out
	if _out == nil {
		consoleOut = false
		_out = &bytes.Buffer{}
	}
	cmd := e.MakeCMD(ctx, _out)
	err := cmd.Run()
	var outRes []byte
	if !consoleOut {
		outRes = _out.(*bytes.Buffer).Bytes()
	}
	if err != nil {
		return false, outRes, err
	}
	return true, outRes, nil
}

func (e *ExecutorImpl) MakeCMD(ctx context.Context, out io.Writer) *exec.Cmd {
	var args []string
	if e.ScriptPath != "" {
		args = append(args, e.ScriptPath)
	}
	args = append(args, e.ScriptArgs...)
	cmd := exec.CommandContext(ctx, e.ShaBang, args...)
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}
	return cmd
}

func (e *ExecutorImpl) ScriptType() ScriptType {
	return e.scriptType
}

// NewPythonExecutor return a python executor
func NewPythonExecutor(scriptPath string, scriptArgs []string) Executor {
	return &ExecutorImpl{
		ShaBang:    "python3",
		ScriptPath: scriptPath,
		ScriptArgs: scriptArgs,
		scriptType: Python,
	}
}

// NewBashExecutor return a bash executor
func NewBashExecutor(scriptPath string, scriptArgs []string) Executor {
	return &ExecutorImpl{
		ShaBang:    "bash",
		ScriptPath: scriptPath,
		ScriptArgs: scriptArgs,
		scriptType: Bash,
	}
}

// NewBinaryExecutor return a binary executor
func NewBinaryExecutor(scriptPath string, scriptArgs []string) Executor {
	return &ExecutorImpl{
		ShaBang:    scriptPath,
		ScriptPath: "",
		ScriptArgs: scriptArgs,
		scriptType: Binary,
	}
}

// NewSQLExecutor return a sql executor
func NewSQLExecutor(scriptPath string) Executor {
	return &SQLExecutor{
		ExecutorImpl: ExecutorImpl{
			ShaBang:    "mysql",
			ScriptPath: scriptPath,
			scriptType: SQL,
		},
	}
}

type SQLExecutor struct {
	Config *viper.Viper
	ExecutorImpl
}

func (e *SQLExecutor) SetConfig(cfg *viper.Viper) {
	e.Config = cfg
}

// MakeCMD make exec.Cmd
// SQLExecutor don't need this method
func (e *SQLExecutor) MakeCMD(ctx context.Context, out io.Writer) *exec.Cmd {
	panic("implement me")
}

func (e *SQLExecutor) Run(ctx context.Context, out io.Writer) (bool, []byte, error) {
	returner := func(err error) (bool, []byte, error) {
		if out != nil {
			if err != nil {
				out.Write([]byte(err.Error()))
				return false, nil, err
			} else {
				fmt.Fprintf(out, "run script: %s success", e.ScriptPath)
				return true, nil, nil
			}
		} else {
			if err != nil {
				return false, []byte(err.Error()), err
			}
		}
		return true, nil, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		e.Config.GetString("MYSQL_USER"),
		e.Config.GetString("MYSQL_PASSWORD"),
		e.Config.GetString("MYSQL_HOST"),
		e.Config.GetInt("MYSQL_PORT"),
		e.Config.GetString("MYSQL_DATABASE"),
	)

	sqlScript, err := os.ReadFile(e.ScriptPath)
	if err != nil {
		return returner(err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return returner(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, string(sqlScript))
	if err != nil {
		return returner(err)
	}
	return returner(nil)
}

func MakeExecutor(filePath string) Executor {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".py":
		return NewPythonExecutor(filePath, nil)
	case ".sh", ".bash":
		return NewBashExecutor(filePath, nil)
	case ".sql":
		return NewSQLExecutor(filePath)
	case "", ".bin":
		return NewBinaryExecutor(filePath, nil)
	default:
		return nil
	}
}
