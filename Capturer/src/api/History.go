package api

import (
	"TianlangCapturer/src/model"
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
)

// Recorder Recorder
type Recorder struct {
	mDataMap       sync.Map
	mSaveFile      string
	mHistoryRecord chan string
}

// RecordHandlerFunc RecordHandlerFunc
type RecordHandlerFunc func(data []byte) (key, value interface{}, err error)

// GetHistory GetHistory
func GetHistory(module IModule, handler RecordHandlerFunc) *Recorder {
	recorder := &Recorder{mHistoryRecord: make(chan string, 128)}
	recorder.mSaveFile = fmt.Sprintf("History.%s.dat", module.Name())
	if recordFile, err := os.Open(recorder.mSaveFile); err == nil {
		defer recordFile.Close()
		reader := bufio.NewReader(recordFile)
		for {
			lineStr, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if data, err := base64.StdEncoding.DecodeString(lineStr); err == nil {
				key, value, err := handler(data)
				if err != nil {
					Error(err.Error())
					continue
				}
				recorder.mDataMap.Store(key, value)
			}
		}
	}
	go recorder.run()
	return recorder
}

// AddRecord AddRecord
func (r *Recorder) AddRecord(key string, value interface{}) {
	r.mDataMap.Store(key, value)
	select {
	case r.mHistoryRecord <- key:
	default:
		Warn("[Recorder]The task queue is full!")
	}
}

// FindRecord FindRecord
func (r *Recorder) FindRecord(key string) (interface{}, bool) {
	if value, ok := r.mDataMap.Load(key); ok {
		return value, true
	}
	return nil, false
}

func (r *Recorder) run() {
	var (
		err  error
		file *os.File
	)
	for key := range r.mHistoryRecord {
		Info("[Doen]Task:%s", key)
		value, _ := r.mDataMap.Load(key)
		data, _ := model.Marshal(value)
		base64Data := base64.StdEncoding.EncodeToString(data)
		if file, err = os.OpenFile(r.mSaveFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
			file.WriteString(base64Data)
			file.WriteString("\r\n")
			file.Close()
		} else {
			Error("[Record]Open '%s' error:%s recordData:\n%s", r.mSaveFile, err.Error(), base64Data)
		}
	}
}
