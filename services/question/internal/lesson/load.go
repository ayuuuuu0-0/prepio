package lesson

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads every world and lesson file under root (content/).
// Unknown fields are rejected so typos in authored files fail loudly.
func Load(root string) (*Content, error) {
	content := &Content{}
	var problems []string

	worldFiles, err := yamlFiles(filepath.Join(root, "worlds"))
	if err != nil {
		return nil, err
	}
	for _, path := range worldFiles {
		var w World
		if err := decodeFile(path, &w); err != nil {
			problems = append(problems, err.Error())
			continue
		}
		content.Worlds = append(content.Worlds, w)
	}

	lessonFiles, err := yamlFiles(filepath.Join(root, "lessons"))
	if err != nil {
		return nil, err
	}
	for _, path := range lessonFiles {
		var l Lesson
		if err := decodeFile(path, &l); err != nil {
			problems = append(problems, err.Error())
			continue
		}
		content.Lessons = append(content.Lessons, l)
	}

	if len(problems) > 0 {
		return nil, errors.New(strings.Join(problems, "\n"))
	}

	sort.SliceStable(content.Worlds, func(i, j int) bool {
		if content.Worlds[i].Order != content.Worlds[j].Order {
			return content.Worlds[i].Order < content.Worlds[j].Order
		}
		return content.Worlds[i].Slug < content.Worlds[j].Slug
	})
	sort.SliceStable(content.Lessons, func(i, j int) bool { return content.Lessons[i].Slug < content.Lessons[j].Slug })
	return content, nil
}

func yamlFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext == ".yaml" || ext == ".yml" {
			files = append(files, path)
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", dir, err)
	}
	sort.Strings(files)
	return files, nil
}

func decodeFile(path string, out any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// Hash returns a stable fingerprint of a lesson's learner-visible content.
// Status is excluded: publishing a lesson does not change what it teaches.
func (l Lesson) Hash() (string, error) {
	raw, err := json.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("hash lesson %s: %w", l.Slug, err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
