package scripts

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const nameSep = "_"

const (
	pre = iota + 1
	run
	post
)

type Script struct {
	Name       string
	SN         int
	IsRollback bool
	Info       os.FileInfo
	Executor
}

type Scripts []Script

// getNameNumber
// 001_init_rollback.sh -> 1
func getNameNumber(name string) int {
	if name == "" {
		return -1
	}
	inxStr := strings.Split(name, nameSep)[0]
	inx, err := cast.ToIntE(inxStr)
	if err != nil {
		return -1
	}
	return inx
}

// isRollback
// 001_init_rollback.sh -> true
// 001_init.sh -> false
func isRollback(name string) bool {
	return strings.HasSuffix(strings.TrimRight(name, filepath.Ext(name)), nameSep+"rollback")
}

// getName
// 001_init_rollback.sh -> init
// 001_init.sh -> init
func getName(name string) string {
	name = strings.TrimRight(name, filepath.Ext(name))
	name = strings.TrimSuffix(name, "_rollback")
	name = strings.Split(name, nameSep)[len(strings.Split(name, nameSep))-1]
	return name
}

// isHide
// _001_init_rollback.sh -> true
func isHide(name string) bool {
	return strings.HasPrefix(name, nameSep)
}

func (s Scripts) GetExecuteScripts() Scripts {
	var scripts Scripts
	for _, script := range s {
		if script.IsRollback {
			continue
		}
		scripts = append(scripts, script)
	}
	return scripts
}

func (s Scripts) GetRollbackScripts(endSN int) Scripts {
	var scripts Scripts
	for _, script := range s {
		if !script.IsRollback {
			continue
		}
		if endSN >= 0 && script.SN > endSN {
			break
		}
		scripts = append(scripts, script)
	}
	return scripts
}

func (s Scripts) HasType(t ScriptType) bool {
	for _, script := range s {
		if script.Executor.ScriptType() == t {
			return true
		}
	}
	return false
}

func (s Scripts) Len() int {
	return len(s)
}

func (s Scripts) Less(i, j int) bool {
	return s[i].SN < s[j].SN
}

func (s Scripts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ScriptManager struct {
	scripts       map[int]Scripts
	BaseDir       string
	Logger        *log.Logger
	Out           io.Writer
	SkipE         bool
	SkipRollbackE bool
	State         map[string]interface{}
	// register check funcs
	checkers []PrepareFunc
	// register get funcs
	getters []PrepareFunc
}

func NewScriptManager(baseDir string, logger *log.Logger, out io.Writer) *ScriptManager {
	return &ScriptManager{
		BaseDir:       baseDir,
		Logger:        logger,
		SkipE:         false,
		Out:           out,
		SkipRollbackE: true,
		scripts:       make(map[int]Scripts),
		State:         make(map[string]interface{}),
		getters: []PrepareFunc{
			MysqlPreGetter,
		},
		checkers: []PrepareFunc{
			MysqlPreChecker,
		},
	}
}

func (m *ScriptManager) makeDirFullPath(dir string) string {
	return filepath.Join(m.BaseDir, dir)
}

func (m *ScriptManager) AddPreScripts(dir string) error {
	scripts, err := Scan(m.makeDirFullPath(dir), true)
	if err != nil {
		return err
	}
	m.scripts[pre] = scripts
	return nil
}

func (m *ScriptManager) AddRunScripts(dir string) error {
	scripts, err := Scan(m.makeDirFullPath(dir), true)
	if err != nil {
		return err
	}
	m.scripts[run] = scripts
	return nil
}

func (m *ScriptManager) AddPostScripts(dir string) error {
	scripts, err := Scan(m.makeDirFullPath(dir), true)
	if err != nil {
		return err
	}
	m.scripts[post] = scripts
	return nil
}

func (m *ScriptManager) runScript(ctx context.Context, stage int, script Script) (*log.Entry, error) {
	logger := m.Logger.WithFields(log.Fields{
		"stage":    stage,
		"rollback": cast.ToString(script.IsRollback),
		"script":   script.Info.Name(),
	})
	logger.Info("run script starting")
	_, out, err := script.Run(ctx, m.Out)
	if err != nil {
		return logger, err
	}
	if len(out) > 0 {
		logger.Info(string(out))
	}
	return logger, nil
}

func (m *ScriptManager) runScripts(ctx context.Context, stage int) error {
	scripts, ok := m.scripts[stage]
	if !ok {
		return nil
	}
	rollback := func(inx int) {
		rollbackScripts := scripts.GetRollbackScripts(inx)
		for i := len(rollbackScripts) - 1; i >= 0; i-- {
			script := rollbackScripts[i]
			logger, err := m.runScript(ctx, stage, script)
			if err != nil && !m.SkipRollbackE {
				logger.WithField("skip_rollback", "false").Info("run rollback script failed, abnormal termination")
				return
			}

		}
	}

	for _, script := range scripts.GetExecuteScripts() {
		logger, err := m.runScript(ctx, stage, script)
		if err != nil && !m.SkipE {
			logger.Info("run script failed, rollback starting")
			rollback(script.SN)
			return err
		}
	}
	return nil
}

func (m *ScriptManager) RunPreScripts(ctx context.Context) error {
	return m.runScripts(ctx, pre)
}

func (m *ScriptManager) RunRunScripts(ctx context.Context) error {
	return m.runScripts(ctx, run)
}

func (m *ScriptManager) RunPostScripts(ctx context.Context) error {
	return m.runScripts(ctx, post)
}

func (m *ScriptManager) Execute(ctx context.Context) error {
	if err := m.RunPreScripts(ctx); err != nil {
		return err
	}
	if err := m.RunRunScripts(ctx); err != nil {
		return err
	}
	if err := m.RunPostScripts(ctx); err != nil {
		return err
	}
	return nil
}

func (m *ScriptManager) GetScripts() map[int]Scripts {
	return m.scripts
}

func (m *ScriptManager) GetStat() map[string]interface{} {
	return m.State
}

func (m *ScriptManager) RunPrepare() error {
	for _, getter := range m.getters {
		if err := getter(m); err != nil {
			return err
		}
	}
	for _, checker := range m.checkers {
		if err := checker(m); err != nil {
			return err
		}
	}
	return nil
}

func Scan(dir string, skipE bool) (Scripts, error) {
	var scripts Scripts
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if !skipE {
				return err
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if isHide(info.Name()) {
			return nil
		}
		executor := MakeExecutor(path)
		if executor == nil {
			return nil
		}
		scripts = append(scripts, Script{
			Name:       getName(info.Name()),
			SN:         getNameNumber(info.Name()),
			IsRollback: isRollback(info.Name()),
			Info:       info,
			Executor:   executor,
		})
		return nil
	})
	sort.Sort(scripts)
	return scripts, err
}
