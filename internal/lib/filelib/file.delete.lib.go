package filelib

import (
	"encoding/json"
	"file-sharing/config"
	"file-sharing/ent"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const DateFormat = "2006-01-02_15-04-05"

type FileInfo struct {
	Name         string `json:"fileName"`
	Types        string `json:"type"`
	Extension    string `json:"extension,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Size         string `json:"size,omitempty"`
}

type DeleteLog struct {
	Large []*FileInfo `json:"largeSize"`
	Small []*FileInfo `json:"smallSize"`
}

func GetType(f os.DirEntry) string {
	if f.IsDir() {
		return "directory"
	}
	return "file"
}

func CreateFileInfo(entries []os.DirEntry) []*FileInfo {
	result := []*FileInfo{}
	for _, f := range entries {
		i, err := f.Info()
		lastModif := ""
		size := ""
		if err == nil {
			lastModif = i.ModTime().Format(DateFormat)
			if !f.IsDir() {
				size = fmt.Sprintf("%.2fMB", float64(i.Size())/float64(config.MB))
			}
		}
		result = append(result, &FileInfo{
			Name:         f.Name(),
			Types:        GetType(f),
			Extension:    GetExtension(f.Name()),
			LastModified: lastModif,
			Size:         size,
		})
	}
	return result
}

func ReadDeleteDir() (large, small []os.DirEntry, err error) {
	err = CreateDir() // prevent error on read
	if err != nil {
		return nil, nil, err
	}

	large, err = os.ReadDir(config.LARGE_PATH)
	if err != nil {
		return nil, nil, err
	}
	small, err = os.ReadDir(config.SMALL_PATH)
	if err != nil {
		return nil, nil, err
	}

	return large, small, nil
}

func CreateDeleteLogFile(large, small []*FileInfo, path string) (*os.File, error) {
	now := time.Now().Format(DateFormat)
	filename := now + ".json"

	os.MkdirAll(path, 0755)
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dl := DeleteLog{Large: large, Small: small}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(dl)

	return file, nil
}

func CreateDeleteLog(files []*ent.File, path string) error {
	large, small, err := ReadDeleteDir()
	if err != nil {
		return err
	}

	ali := CreateFileInfo(large)
	asi := CreateFileInfo(small)

	li := []*FileInfo{}
	si := []*FileInfo{}

	for _, file := range files {
		if GetPathBySize(file.FileSize) == config.LARGE_PATH {
			for _, i := range ali {
				if i.Name == file.FileName {
					li = append(li, i)
					break
				}
			}
		} else {
			for _, i := range asi {
				if i.Name == file.FileName {
					si = append(si, i)
					break
				}
			}
		}
	}

	if _, err := CreateDeleteLogFile(li, si, path); err != nil {
		return err
	}

	return nil
}
