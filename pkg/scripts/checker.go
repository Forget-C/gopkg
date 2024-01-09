package scripts

import (
	"fmt"
	"github.com/spf13/viper"
)

type Manager interface {
	GetStat() map[string]interface{}
	GetScripts() map[int]Scripts
}

type PrepareFunc func(m Manager) error

func hasScript(m Manager, t ScriptType) bool {
	has := false
	for _, ss := range m.GetScripts() {
		if ss.HasType(t) {
			has = true
			break
		}
	}
	return has
}

func MysqlPreGetter(m Manager) error {
	if !hasScript(m, SQL) {
		return nil
	}
	config := viper.New()
	config.AutomaticEnv()
	config.SetDefault("MYSQL_PORT", 3306)
	config.SetDefault("MYSQL_USER", "root")
	stat := m.GetStat()
	stat["SQLConfig"] = config
	return nil
}

func MysqlPreChecker(m Manager) error {
	if !hasScript(m, SQL) {
		return nil
	}
	stat := m.GetStat()
	config, ok := stat["SQLConfig"]
	if !ok {
		return fmt.Errorf("mysql config not found")
	}
	vCfg, ok := config.(*viper.Viper)
	if !ok {
		return fmt.Errorf("mysql config type error")
	}
	if vCfg.GetInt("MYSQL_PORT") == 0 || vCfg.GetString("MYSQL_USER") == "" ||
		vCfg.GetString("MYSQL_PASSWORD") == "" || vCfg.GetString("MYSQL_HOST") == "" ||
		vCfg.GetString("MYSQL_DATABASE") == "" {
		return fmt.Errorf("mysql config error")
	}
	for _, ss := range m.GetScripts() {
		for _, s := range ss {
			if s.ScriptType() == SQL {
				s.Executor.(*SQLExecutor).SetConfig(vCfg)
			}
		}
	}
	return nil
}
