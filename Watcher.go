package AnimaKit

import (
	"sync"
	"time"

	"github.com/radovskyb/watcher"
)

var Files2WatchLock *sync.RWMutex
var Files2Watch = make(map[string]bool)
var ChangedSignal = make(chan bool)
var StopSignal = make(chan bool)

func init() {
	Files2WatchLock = new(sync.RWMutex)
}

func AddFileToWatch(path string) {
	Files2WatchLock.Lock()
	Files2Watch[path] = true
	Files2WatchLock.Unlock()
}

func WatchFiles() {
	// creates a new file watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				TheLog.NoticeF("File changed: %s", event.Path)
				ChangedSignal <- true
				w.Close()
				return
			case err := <-w.Error:
				TheLog.Error(err)
			case <-w.Closed:
				return
			}
		}
	}()

	Files2WatchLock.RLock()
	defer Files2WatchLock.RUnlock()
	for path, _ := range Files2Watch {
		if err := w.Add(path); err != nil {
			TheLog.Error(err)
		}
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		TheLog.Error(err)
	}
}
