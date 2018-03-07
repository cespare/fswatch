// fswatch is a small Go library that builds on top of
// github.com/fsnotify/fsnotify and adds a few commonly-needed extra features,
// including recursive watches.
package fswatch

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
)

const enableLog = false

type watcher struct {
	dir              string
	coalesceInterval time.Duration
	w                *fsnotify.Watcher
	events           chan []string
	errors           chan error
}

// Watch creates a recursive file watch for dir.
//
// It returns two channels; they have meanings similar to fsnotify.Watcher,
// except that the events channel returns slices of modified files/directories).
//
// The watcher coalesces quick sequences of events into a single event slice,
// using a time window specified by coalesceInterval.
func Watch(dir string, coalesceInterval time.Duration) (events chan []string, errs chan error, err error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	w := &watcher{
		dir:              dir,
		coalesceInterval: coalesceInterval,
		w:                fw,
		events:           make(chan []string),
		errors:           make(chan error),
	}
	if err := w.addDirRecursive(dir); err != nil {
		return nil, nil, err
	}
	go w.filter()
	go func() {
		for err := range fw.Errors {
			w.errors <- err
		}
	}()
	return w.events, w.errors, nil
}

const chmodMask fsnotify.Op = ^fsnotify.Op(0) ^ fsnotify.Chmod

func (w *watcher) filter() {
	timer := time.NewTimer(0)
	<-timer.C
	defer timer.Stop()
	timerStarted := false
	seen := make(map[string]struct{})
	for {
		select {
		case ev, ok := <-w.w.Events:
			if !ok {
				return
			}
			// Ignore events that are *only* CHMOD to work around Spotlight.
			if ev.Op&chmodMask == 0 {
				continue
			}
			if !timerStarted {
				timer.Reset(w.coalesceInterval)
				timerStarted = true
			}
			seen[ev.Name] = struct{}{}
			if ev.Op&fsnotify.Create != 0 && isDir(ev.Name) {
				if err := w.addDirRecursive(ev.Name); err != nil {
					w.errors <- err
				}
			}
		case <-timer.C:
			var names []string
			for name := range seen {
				names = append(names, name)
			}
			sort.Strings(names)
			w.events <- names
			seen = make(map[string]struct{})
			timerStarted = false
		}
	}
}

func (w *watcher) addDirRecursive(dir string) error {
	return filepath.Walk(dir, w.addDirsWalkFunc())
}

func (w *watcher) addDirsWalkFunc() filepath.WalkFunc {
	return func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if enableLog {
			log.Println("Adding watch for", name)
		}
		return w.w.Add(name)
	}
}

func isDir(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
