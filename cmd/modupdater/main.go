package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/kroyser123/go-mod-updater/internal/git"
	"github.com/kroyser123/go-mod-updater/internal/logger"
	"github.com/kroyser123/go-mod-updater/internal/modfinder"
	"github.com/kroyser123/go-mod-updater/internal/modparser"
	"github.com/kroyser123/go-mod-updater/internal/output"
	"github.com/kroyser123/go-mod-updater/internal/version"
)

func main() {
	repo := flag.String("repo", "", "Git repository URL")
	token := flag.String("token", "", "Access token")
	showAll := flag.Bool("all", false, "Show all dependencies, not only outdated")
	jsonOut := flag.Bool("json", false, "Output results in JSON format")
	debug := flag.Bool("debug", false, "Enable debug logs")

	flag.Parse()

	log := logger.NewLogger(*debug)

	if *repo == "" {
		log.Error("Missing required flag: -repo")
		return
	}

	// клонируем репозиторий

	tmpDir, err := git.Clone(*repo, *token, log)
	if err != nil {
		log.Error("Clone failed: %v", err)
		return
	}

	// Ищем go.mod файлы

	modFiles, err := modfinder.Find(tmpDir, log)
	if err != nil {
		log.Error("Failed to find go.mod: %v", err)
		return
	}
	var results []output.ModuleResult

	// Проходим по каждому модулю и чекаем апдейты

	for _, modPath := range modFiles {
		mod, err := modparser.Parse(modPath, log)
		if err != nil {
			results = append(results, output.NewModuleResult(modPath, "", nil, err))
			continue
		}

		workDir := filepath.Dir(modPath)
		statuses, err := version.Check(mod, *showAll, log, workDir)

		results = append(results,
			output.NewModuleResult(modPath, mod.ModulePath, statuses, err),
		)
	}

	// json вывод
	if *jsonOut {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	// CLI вывод
	output.Print(results)
}
