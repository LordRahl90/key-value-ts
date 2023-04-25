package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"key-value-ts/domains/entities"
	"key-value-ts/domains/storage"
)

const ext = ".csv"

var (
	_ storage.Storer = (*FileStorer)(nil)
)

// FileStorer storer implementation to save records to files
type FileStorer struct {
	path  string
	cache map[string]entities.Sequence
	m     sync.Mutex
}

// New returns a FileStorer implementation
func New(path string) (storage.Storer, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if _, err := os.Stat(path); err != nil {
		if inErr := os.MkdirAll(filepath.Dir(path), 0700); inErr != nil {
			if !os.IsExist(inErr) {
				return nil, err
			}
		}
	}
	return &FileStorer{
		path:  path,
		cache: make(map[string]entities.Sequence),
	}, nil
}

// Get retrieves the sequence from the map first, then checks the file.
func (fs *FileStorer) Get(key string, timestamp int64) (string, error) {
	if v, ok := fs.fromCache(key, timestamp); ok {
		return v.Value, nil
	}
	path := fmt.Sprintf("%s/%s%s", strings.TrimSuffix(fs.path, "/"), key, ext)
	csvFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	result := ""
	r := csv.NewReader(csvFile)
	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return result, nil
			}
			return result, err
		}
		if len(rec) < 3 {
			continue
		}
		ts, err := strconv.Atoi(rec[1])
		if err != nil {
			return result, err
		}

		if timestamp == int64(ts) {
			fs.saveToCache(entities.Sequence{
				Key:       key,
				Timestamp: int64(ts),
				Value:     rec[2],
			})

			result = rec[2]
		}
	}
}

// Save stores the sequence into files
func (fs *FileStorer) Save(e entities.Sequence) error {
	fileName := fmt.Sprintf("%s/%s%s", strings.TrimSuffix(fs.path, "/"), e.Key, ext)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s,%d,%s\n", e.Key, e.Timestamp, e.Value)); err != nil {
		return err
	}
	fs.saveToCache(e)
	return nil
}

func (fs *FileStorer) saveToCache(e entities.Sequence) {
	fs.m.Lock()
	defer fs.m.Unlock()
	fs.cache[mapKey(e.Key, e.Timestamp)] = e
}

func (fs *FileStorer) fromCache(key string, timestamp int64) (*entities.Sequence, bool) {
	val, ok := fs.cache[mapKey(key, timestamp)]
	return &val, ok
}

func mapKey(key string, timestamp int64) string {
	return fmt.Sprintf("%s_%d", key, timestamp)
}
