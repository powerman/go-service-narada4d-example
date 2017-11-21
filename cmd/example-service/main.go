package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/powerman/go-service-narada4d-example/api"
	"github.com/powerman/go-service-narada4d-example/dal"
	"github.com/powerman/structlog"
)

var (
	app = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	ver string
	log = structlog.New()
	cfg struct {
		version  bool
		logLevel string
		apiKey   string
	}
)

func init() {
	flag.BoolVar(&cfg.version, "version", false, "print version")
	flag.StringVar(&cfg.logLevel, "log.level", "debug", "log level (debug|info|warn|err)")
	flag.StringVar(&cfg.apiKey, "api.key", os.Getenv("API_KEY"), "secret required to access API ($API_KEY)")
}

// Init provides common initialization for both app and tests.
func Init() {
	structlog.DefaultLogger.
		SetSuffixKeys(
			structlog.KeyStack,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyUnit: " %6[2]s:",
		})
	time.Local = time.UTC
}

func main() {
	Init()
	flag.Parse()

	if cfg.version {
		fmt.Println(app, ver, runtime.Version())
		os.Exit(0)
	}

	structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(cfg.logLevel))
	structlog.DefaultLogger.SetDefaultKeyvals(structlog.KeyUnit, "main")

	log.Info("started")
	defer log.Info("finished")

	if err := dal.Init(); err != nil {
		log.Fatal(err)
	}

	if err := api.Init(cfg.apiKey, ver); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Minute)
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		panic("oops")
	}
}
