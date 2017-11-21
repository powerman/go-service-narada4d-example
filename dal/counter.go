package dal

import (
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/powerman/go-service-narada4d-example/narada4d"
)

var (
	counterMutex   sync.Mutex
	counterPath    string
	counterPathTmp string
)

// Count will increment and return value of db/counter.
func Count() (int, error) {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	narada4d.SharedLock()
	defer narada4d.Unlock()

	buf, err := ioutil.ReadFile(counterPath)
	if err != nil {
		return 0, err
	}
	counter, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}
	counter++
	err = ioutil.WriteFile(counterPathTmp, []byte(strconv.Itoa(counter)), 0666)
	if err != nil {
		return 0, err
	}
	err = os.Rename(counterPathTmp, counterPath)
	if err != nil {
		return 0, err
	}
	return counter, nil
}
