package dal

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/powerman/go-service-narada4d-example/narada4d"
	. "gopkg.in/check.v1"
)

type TestSchema struct{}

var _ = Suite(&TestSchema{})

func (s *TestSchema) TestRight(c *C) {
	// crasher(c, 0, "", func() {
	setc := setSchemaVersion()
	donec := make(chan struct{})
	setc <- narada4d.SetVersion{Version: schemaVersion, Done: donec}
	<-donec
	// })
}

func (s *TestSchema) TestWrong(c *C) {
	crasher(c, 1, ".*incompatible schema version.*\n", func() {
		setc := setSchemaVersion()
		donec := make(chan struct{})
		setc <- narada4d.SetVersion{Version: "0", Done: donec}
		<-donec
	})
}

func crasher(c *C, exitCode int, stderrRegex string, test func()) {
	if os.Getenv("BE_CRASHER") == "1" {
		test()
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	caller := parts[len(parts)-1]
	cmd := exec.Command(os.Args[0], "-check.f=^"+caller+"$")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	_, err := cmd.Output()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		c.Check(e.Sys().(syscall.WaitStatus).ExitStatus(), Equals, exitCode)
		c.Check(string(e.Stderr), Matches, stderrRegex)
	} else {
		c.Errorf("process ran with err %v, want exit status 1", err)
	}
}
