package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	dirName  = ".aitutor"
	fileName = "progress.json"
)

// ProgressData is the JSON-serialized progress state.
type ProgressData struct {
	CompletedLessons map[int]bool `json:"completed_lessons"`
	LastLessonIdx    int          `json:"last_lesson_idx"`
}

func progressPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

// Load reads progress from disk. Returns empty data if file doesn't exist.
func Load() (ProgressData, error) {
	data := ProgressData{CompletedLessons: make(map[int]bool)}

	path, err := progressPath()
	if err != nil {
		return data, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return ProgressData{CompletedLessons: make(map[int]bool)}, nil
	}
	if data.CompletedLessons == nil {
		data.CompletedLessons = make(map[int]bool)
	}
	return data, nil
}

// Save writes progress to disk.
func Save(data ProgressData) error {
	path, err := progressPath()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
