package explorer

import (
	"log"
	"os"

	"github.com/ShakedGold/Gole/pkg/widgets/entry"
	"github.com/fsnotify/fsnotify"
)

func Home() (*entry.Entries, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &entry.Entries{}, err
	}
	return entry.ReadPath(homeDir)
}

// Watcher watches for changes in the filesystem
func Watcher(path string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

func Watch(watcher *fsnotify.Watcher, entries *entry.Entries) {
	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}

			// if there is a change update the entries
			entries, _ = entry.ReadPath(entries.Path)
			newEntries, err := entries.Prepare()
			if err != nil {
				log.Println(err)
			}
			entries = newEntries
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("Error:", err)
		}
	}
}
