package modfinder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/kroyser123/go-mod-updater/internal/logger"
)

func Find(root string, log *logger.Logger) ([]string, error) {

	// Ищем go.mod

	log.Debug("Searching for go.mod in: %s", root)
	var result []string

	// если root пустой - ошибка

	if root == "" {
		log.Error("empty root: %v", root)
		return result, fmt.Errorf("Root is empty")
	}
	info, err := os.Stat(root)
	if err != nil {
		log.Error("cannot access root: %v", err)
		return nil, fmt.Errorf("cannot access %s: %w", root, err)
	}
	if !info.IsDir() {
		log.Error("root is not a directory: %s", root)
		return nil, fmt.Errorf("root is not a directory: %s", root)
	}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		// проходим по директории, скипаем .git
		// ошибка доступа к файлу - скип

		if err != nil {
			log.Debug("Skipping %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		log.Error("Error walking directory: %v", err)
		return nil, err
	}
	sort.Strings(result)
	return result, nil
}
