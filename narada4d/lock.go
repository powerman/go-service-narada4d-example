package narada4d

import (
	"errors"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/powerman/structlog"
)

var log = structlog.New().SetDefaultKeyvals(structlog.KeyUnit, "narada")

// Errors returned by Init.
var (
	ErrDirNotSet         = errors.New("NARADA4D_DIR is not set")
	ErrDirNotDir         = errors.New("NARADA4D_DIR is not a directory")
	ErrDirNotInitialized = errors.New("NARADA4D_DIR is not initialized")
)

const (
	versionFileName   = ".version"
	lockFileName      = ".lock"
	lockQueueFileName = ".lock.queue"
)

var (
	// Dir is $NARADA4D_DIR.
	Dir           string
	setVersion    chan<- SetVersion
	versionPath   string
	skipLock      bool
	lockFile      *os.File
	lockQueueFile *os.File
	lockFD        int
	lockQueueFD   int
	isLockedMutex sync.Mutex
	isLocked      bool
)

// SetVersion is an event sent after successful locking.
// Event consumer should check is it compatible with Version and either
// setup data access accordingly to Version and then send to Done or
// shutdown application without sending to Done.
type SetVersion struct {
	Version string
	Done    chan<- struct{}
}

// Init must be called before using this package.
func Init(c chan<- SetVersion) error {
	setVersion = c

	Dir = os.Getenv("NARADA4D_DIR")
	if Dir == "" {
		return log.Err(ErrDirNotSet)
	}
	Dir = path.Clean(Dir)
	fi, err := os.Stat(Dir)
	if err != nil || !fi.IsDir() {
		return log.Err(ErrDirNotDir)
	}

	versionPath = path.Join(Dir, versionFileName)
	fi, err = os.Lstat(versionPath)
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		return log.Err(ErrDirNotInitialized)
	}

	skipLock = os.Getenv("NARADA4D_SKIP_LOCK") != ""

	if !skipLock {
		lockFile, err = os.Open(path.Join(Dir, lockFileName))
		if err != nil {
			return log.Err(err)
		}
		lockFD = int(lockFile.Fd())

		lockQueueFile, err = os.Open(path.Join(Dir, lockQueueFileName))
		if err != nil {
			return log.Err(err)
		}
		lockQueueFD = int(lockQueueFile.Fd())
	}

	return nil
}

// ExclusiveLock will set exclusive lock, send SetVersion event and wait
// for the response.
// You must Unlock before calling ExclusiveLock or SharedLock again.
func ExclusiveLock() {
	lock(syscall.LOCK_EX)
}

// SharedLock will set shared lock, send SetVersion event and wait
// for the response.
// You must Unlock before calling ExclusiveLock or SharedLock again.
func SharedLock() {
	lock(syscall.LOCK_SH)
}

// Unlock will release lock set by ExclusiveLock or SharedLock.
// You must call ExclusiveLock or SharedLock before Unlock.
func Unlock() {
	isLockedMutex.Lock()
	defer isLockedMutex.Unlock()
	if !isLocked {
		panic("not locked")
	}
	isLocked = false

	if err := syscall.Flock(lockFD, syscall.LOCK_UN); err != nil {
		panic(err)
	}
}

func lock(how int) {
	isLockedMutex.Lock()
	defer isLockedMutex.Unlock()
	if isLocked {
		panic("already locked")
	}
	isLocked = true

	if !skipLock {
		if err := syscall.Flock(lockQueueFD, syscall.LOCK_EX); err != nil {
			panic(err)
		}
		if err := syscall.Flock(lockFD, how); err != nil {
			panic(err)
		}
		if err := syscall.Flock(lockQueueFD, syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}

	version, err := os.Readlink(versionPath)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	setVersion <- SetVersion{
		Version: version,
		Done:    done,
	}
	<-done
}
