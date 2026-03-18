package modparser

import (
	"fmt"
	"os"

	"github.com/kroyser123/go-mod-updater/internal/logger"
	"golang.org/x/mod/modfile"
)

// создаем структуру require, которая представляет собой одну зависимость из файла go.mod

type Require struct {
	Path     string // путь зависимости
	Version  string
	Indirect bool
}

// создаем структуру ModFile, которая представляет собой весь файл go.mod

type ModFile struct {
	ModulePath string // путь текущего модуля
	Version    string // GO версия
	Requires   []Require
}

// читаем и парсим файл go.mod

func Parse(modPath string, log *logger.Logger) (*ModFile, error) {
	log.Debug("Parsing go.mod: %s", modPath)
	data, err := os.ReadFile(modPath)
	if err != nil {
		log.Error("Failed to read mod file: %v", err)
		return nil, fmt.Errorf("Failed to read mod file %v", err)
	}

	// парсим содержимое

	file, err := modfile.Parse(modPath, data, nil)
	if err != nil {
		log.Error("Failed to parse: %v", err)
		return nil, fmt.Errorf("Fail to parse the file %v", err)
	}

	// результат будет храниться в res
	// выделяем память ровно под столько элементов, сколько в модуле

	res := &ModFile{
		Requires: make([]Require, 0, len(file.Require)),
	}

	// если модуль пуст возвращаем ошибку, если версия пустая - нормально

	if file.Module != nil {
		res.ModulePath = file.Module.Mod.Path
	} else {
		log.Error("ModulePath is empty")
		return nil, fmt.Errorf("ModulePath is empty")
	}
	if file.Go != nil {
		res.Version = file.Go.Version
	}

	// заполняем Requires

	for _, r := range file.Require {
		if r.Mod.Path == "" {
			log.Debug("skipping empty require")
			continue
		}

		log.Debug("Adding require: %s %s", r.Mod.Path, r.Mod.Version)

		res.Requires = append(res.Requires, Require{
			Path:     r.Mod.Path,
			Version:  r.Mod.Version,
			Indirect: r.Indirect,
		})
	}
	log.Debug("%d dependencies was/were found", len(res.Requires))
	return res, nil
}
