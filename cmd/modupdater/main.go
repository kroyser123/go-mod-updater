package main

import (
	"flag"
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
	flag.Parse()

	log := logger.NewLogger(true)

	if *repo == "" {
		log.Error("No repo URL provided")
		return
	}

	// 1. Clone
	tmpDir, err := git.Clone(*repo, *token, log)
	if err != nil {
		log.Error("Clone failed: %v", err)
		return
	}

	// 2. Find go.mod
	modFiles, err := modfinder.Find(tmpDir, log)
	if err != nil {
		log.Error("Find failed: %v", err)
		return
	}

	var results []output.ModuleResult

	// 3. Process each go.mod
	for _, modPath := range modFiles {
		mod, err := modparser.Parse(modPath, log)
		if err != nil {
			results = append(results, output.NewModuleResult(modPath, "", nil, err))
			continue
		}

		workDir := filepath.Dir(modPath)

		statuses, err := version.Check(mod, false, log, workDir)

		results = append(results,
			output.NewModuleResult(modPath, mod.ModulePath, statuses, err),
		)
	}

	// 4. Print
	output.Print(results)
}
