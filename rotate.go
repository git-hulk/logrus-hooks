package hooks

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	RotateByDay = iota
	RotateByHour
)

const (
	defaultDayTimePattern  = "20060102"
	defaultHourTimePattern = "20060102-15"
)

type RotateHook struct {
	dir          string
	name         string
	currFileTime string
	writer       *os.File
	timePattern  string
	lock         sync.Mutex
	logger       *logrus.Logger
}

func NewRotateHook(logger *logrus.Logger, dir, name string) (*RotateHook, error) {
	rh := new(RotateHook)
	rh.name = name
	rh.timePattern = defaultDayTimePattern
	rh.logger = logger
	rh.dir = dir

	writer, err := rh.openNewFile()
	if err != nil {
		return nil, err
	}
	rh.writer = writer
	logger.Out = writer
	return rh, nil
}

func (rh *RotateHook) openNewFile() (*os.File, error) {
	_, err := os.Stat(rh.dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rh.dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	newFileTime := time.Now().Format(rh.timePattern)
	newFilename := fmt.Sprintf("%s/%s.log.%s", rh.dir, rh.name, newFileTime)
	newWriter, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	rh.currFileTime = newFileTime

	return newWriter, nil
}

func ReleaseRotateHook(rh *RotateHook) {
	rh.writer.Close()
}

func (rh *RotateHook) needRotate() bool {
	return rh.currFileTime != time.Now().Format(rh.timePattern)
}

func (rh *RotateHook) rotate() error {
	rh.lock.Lock()
	defer rh.lock.Unlock()

	// check again, someone may have rotated at the same time
	if !rh.needRotate() {
		return nil
	}

	oldWriter := rh.writer
	newWriter, err := rh.openNewFile()
	if err != nil {
		return err
	}
	rh.writer = newWriter
	rh.logger.Out = newWriter
	err = oldWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rh *RotateHook) SetRotateType(rtype int) {
	switch rtype {
	case RotateByDay:
		rh.timePattern = defaultDayTimePattern
	case RotateByHour:
		rh.timePattern = defaultHourTimePattern
	}
}

func (rh *RotateHook) Fire(entry *logrus.Entry) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if rh.needRotate() {
		return rh.rotate()
	}
	return nil
}

func (rh *RotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
