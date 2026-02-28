package ant

import (
	"github.com/yuancore/go-zen/db/adb"
	"github.com/yuancore/go-zen/os/config"
	"gorm.io/gorm"
)

// Db Get database connection
func Db(name ...string) *gorm.DB {
	key := ""

	if len(name) > 0 {
		key = name[0]
	}

	val, ok := adb.Master[key]
	if ok {
		return val
	} else {
		key = config.GetString("connections.0.name")
	}

	return adb.Master[key]
}

// CloseDb Close database connection
func CloseDb() {
	adb.Close()
}
