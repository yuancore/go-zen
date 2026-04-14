package serve

import (
	"github.com/yuancore/go-zen/adapter/db/zdb"
	"github.com/yuancore/go-zen/zen"
)

// 定义与配置文件匹配的中间结构体
type rawConnection struct {
	Name            string `mapstructure:"name"`
	Type            string `mapstructure:"type"`
	Hostname        string `mapstructure:"hostname"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Params          string `mapstructure:"params"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

// loadZdb 从顶层 [[connections]] 配置创建 zdb 组件
func loadZdb(cfg zen.Config) (*zdb.Component, error) {
	if !cfg.IsSet("connections") {
		return nil, nil
	}

	var rawConns []rawConnection
	if err := cfg.Unmarshal("connections", &rawConns); err != nil {
		return nil, err
	}
	if len(rawConns) == 0 {
		return nil, nil
	}

	settings := zdb.Settings{}
	for _, rc := range rawConns {
		conn := zdb.ConnectionConfig{
			Name:            rc.Name,
			Driver:          rc.Type,
			Host:            rc.Hostname,
			Port:            rc.Port,
			Username:        rc.Username,
			Password:        rc.Password,
			Database:        rc.Database,
			Parameters:      rc.Params,
			MaxIdleConns:    rc.MaxIdleConns,
			MaxOpenConns:    rc.MaxOpenConns,
			ConnMaxLifetime: rc.ConnMaxLifetime,
			ConnMaxIdleTime: rc.ConnMaxIdleTime,
		}
		settings.Connections = append(settings.Connections, conn)
		if settings.Default == "" {
			settings.Default = conn.Name
		}
	}
	return zdb.New(zdb.WithSettings(settings)), nil
}
