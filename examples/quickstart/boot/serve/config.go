package serve

import (
	"fmt"

	zconfig "github.com/yuancore/go-zen/adapter/config/viper"
	"github.com/yuancore/go-zen/zen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// loadConfig 加载配置文件，不做任何魔法路径处理
func loadConfig(path string) (zen.Config, error) {
	cfg := zconfig.NewConfig()
	if err := cfg.Load(path); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}
	return cfg, nil
}

// newDB 根据配置初始化数据库连接（当前仅支持 SQLite，可扩展）
func newDB(cfg zen.Config) (*gorm.DB, error) {
	driver := cfg.GetString("database.driver")
	dsn := cfg.GetString("database.dsn")

	switch driver {
	case "sqlite", "sqlite3":
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}
