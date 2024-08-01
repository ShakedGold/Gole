package explorer

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"gioui.org/widget"
	"github.com/fsnotify/fsnotify"
	"github.com/skratchdot/open-golang/open"
)

type Entry struct {
	Path      string
	Alias     string
	IsFolder  bool
	Clickable *widget.Clickable
}

type Entries struct {
	Entries []Entry
	Path    string
}

func Home() (Entries, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Entries{}, err
	}
	return ReadPath(homeDir)
}

func ReadPath(path string) (Entries, error) {
	osEntries, err := os.ReadDir(path)

	if err != nil {
		return Entries{}, err
	}

	entries := Entries{
		Entries: []Entry{},
		Path:    path,
	}

	for _, entry := range osEntries {
		entries.Entries = append(entries.Entries, Entry{Path: filepath.Join(path, entry.Name()), IsFolder: entry.IsDir(), Clickable: new(widget.Clickable)})
	}

	return entries, nil
}

func (entries Entries) Prepare() Entries {
	// add .. and . entry
	entries.Entries = append([]Entry{{Path: filepath.Join(entries.Path, ".."), Alias: "..", IsFolder: true}}, entries.Entries...)
	entries.Entries = append(entries.Entries, Entry{Path: entries.Path, Alias: ".", IsFolder: true})

	// sort entries
	sort.Slice(entries.Entries, ByIsDir(entries.Entries))

	return entries
}

func ByIsDir(entries []Entry) func(a, b int) bool {
	return func(a, b int) bool {
		// sort by IsDir, then by Name
		if entries[a].IsFolder == entries[b].IsFolder {
			return entries[a].Path < entries[b].Path
		}
		return entries[a].IsFolder
	}
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

func Watch(watcher *fsnotify.Watcher, entries *Entries) {
	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}

			// if there is a change update the entries
			*entries, _ = ReadPath(entries.Path)
			*entries = entries.Prepare()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("Error:", err)
		}
	}
}

func (e Entry) EntryAction(watcher *fsnotify.Watcher) (Entries, error) {
	if e.IsFolder {
		entries, err := ReadPath(e.Path)
		if err != nil {
			log.Println(err)
			return Entries{}, err
		}

		entries = entries.Prepare()

		// remove from watch
		watcher.Remove(e.Path)

		// watch new folder
		watcher.Add(entries.Path)

		return entries, nil
	} else {
		open.Run(e.Path)
	}

	return Entries{}, nil
}
