package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"

	"github.com/kroyser123/go-mod-updater/internal/logger"
	"github.com/kroyser123/go-mod-updater/internal/modparser"
	"golang.org/x/mod/semver"
)

// типы апдейтов
type Update string

const (
	NoUpdate Update = "none"
	Patch    Update = "patch"
	Minor    Update = "minor"
	Major    Update = "major"
)

// структура по каждой зависимости

type Dependecies struct {
	Path       string
	Current    string
	Latest     string
	Indirect   bool
	NeedUpdate bool
	UpdateType Update
	Error      error
}

// структура под go list -m -u -json all

type GoList struct {
	Path     string `json:"Path"`
	Version  string `json:"Version"`
	Indirect bool   `json:"Indirect,omitempty"`
	Update   *struct {
		Version string `json:"Version"`
	} `json:"Update,omitempty"`
	Error *struct {
		Err error `json:"Err"`
	} `json:"Error,omitempty"`
}

// проверяем зависимости функцией Check

func Check(mod *modparser.ModFile, TurnOnIndirect bool, log *logger.Logger, WorkDir string) ([]Dependecies, error) {

	// запускаем go list

	cmd := exec.Command("go", "list", "-m", "-u", "-json", "all")
	cmd.Dir = WorkDir
	Output, err := cmd.Output()
	if err != nil {
		log.Error("failed to exec go list: %v", err)
		return nil, fmt.Errorf("failed to exec go list: %v", err)
	}
	modules := make(map[string]GoList)

	// парсим json поток

	decoder := json.NewDecoder(bytes.NewReader(Output))

	for {
		var g GoList
		if err := decoder.Decode(&g); err != nil {
			break
		}
		modules[g.Path] = g
	}
	log.Debug("Parsed modules from go list: %d", len(modules))

	// обработка зависимостей
	// добавлена параллельная обработка

	var wg sync.WaitGroup
	var mu sync.Mutex
	statuses := make([]Dependecies, 0, len(mod.Requires))
	for _, req := range mod.Requires {
		if !TurnOnIndirect && req.Indirect {
			continue
		}
		wg.Add(1)
		go func(req modparser.Require) {
			defer wg.Done()
			st := Dependecies{
				Path:     req.Path,
				Current:  req.Version,
				Indirect: req.Indirect,
			}

			// проверяем существует ли в modules распаршенный модуль из modfile

			m, ok := modules[req.Path]
			if !ok {
				st.Error = fmt.Errorf("module was not found in go list")
			}
			if m.Error != nil {
				st.Error = m.Error.Err
			}

			// есть или нет обновления

			if m.Update == nil {
				st.Latest = m.Version
				st.UpdateType = NoUpdate
				st.NeedUpdate = false
			} else {
				st.Latest = m.Update.Version
				st.UpdateType = UpdateType(req.Version, m.Update.Version)
				st.NeedUpdate = true
			}
			mu.Lock()
			statuses = append(statuses, st)
			mu.Unlock()
		}(req)
	}
	wg.Wait()
	return statuses, nil
}

// определяем уровень update'а с помощью официальной библиотеки semver

func UpdateType(current, latest string) Update {
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return NoUpdate
	}
	if semver.Major(current) != semver.Major(latest) {
		return Major
	}
	if semver.MajorMinor(current) != semver.MajorMinor(latest) {
		return Minor
	}
	if semver.Compare(current, latest) < 0 {
		return Patch
	}
	return NoUpdate
}
