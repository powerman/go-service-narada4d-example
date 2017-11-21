package dal

import (
	"os"

	"github.com/powerman/go-service-narada4d-example/narada4d"
)

const schemaVersion = "1"

func setSchemaVersion() chan<- narada4d.SetVersion {
	c := make(chan narada4d.SetVersion)
	go func() {
		v := <-c
		for ; v.Version == schemaVersion; v = <-c {
			v.Done <- struct{}{}
		}
		log.Err("incompatible schema version", "required", schemaVersion, "current", v.Version)
		os.Exit(1)
	}()
	return c
}
