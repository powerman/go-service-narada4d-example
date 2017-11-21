package dal

import (
	"path"

	"github.com/powerman/go-service-narada4d-example/narada4d"
	"github.com/powerman/structlog"
)

var log = structlog.New()

// Init must be called before using this package.
func Init() error {
	if err := narada4d.Init(setSchemaVersion()); err != nil {
		return err
	}

	counterPath = path.Join(narada4d.Dir, "db/counter")
	counterPathTmp = counterPath + ".tmp"
	return nil
}
