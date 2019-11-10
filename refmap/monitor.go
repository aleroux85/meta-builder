package refmap

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	Watcher *fsnotify.Watcher
	Stopped chan bool
	Error   *error
}

func NewMonitor() *Monitor {
	m := new(Monitor)
	var temp error
	m.Error = &temp
	return m
}

func (m *Monitor) SetWatcher() error {
	if *m.Error != nil {
		return nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		*m.Error = err
		return err
	}
	m.Watcher = w

	m.Stopped = make(chan bool)
	return nil
}

func (m *Monitor) Close() {
	err := m.Watcher.Close()
	if err != nil {
		*m.Error = err
	}
	<-m.Stopped
}

func (m *Monitor) AddDirectories(s string, ignore map[string]struct{}) error {
	if *m.Error != nil {
		return nil
	}

	if ignore == nil {
		ignore = make(map[string]struct{})
	}

	info, err := os.Stat(s)
	if err != nil {
		*m.Error = err
		return err
	}

	if info.IsDir() {
		err := m.Watcher.Add(s)
		if err != nil {
			*m.Error = err
			return err
		}

		file, err := os.Open(s)
		if err != nil {
			*m.Error = err
			return err
		}

		files, err := file.Readdir(0)
		if err != nil {
			*m.Error = err
			return err
		}

		for _, iFile := range files {
			if iFile.IsDir() {
				_, found := ignore[filepath.Join(s, iFile.Name())]
				if !found {
					m.AddDirectories(filepath.Join(s, iFile.Name()), ignore)
				}
			}
		}
	}

	return nil
}
