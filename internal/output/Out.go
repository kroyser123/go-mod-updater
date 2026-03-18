package output

import (
	"fmt"

	"github.com/kroyser123/go-mod-updater/internal/version"
)

type ModuleResult struct {
	ModFilePath string
	ModulePath  string
	Statuses    []version.Dependecies
	Err         error
}

// Конструктор — создаёт готовый объект
func NewModuleResult(path string, module string, statuses []version.Dependecies, err error) ModuleResult {
	return ModuleResult{
		ModFilePath: path,
		ModulePath:  module,
		Statuses:    statuses,
		Err:         err,
	}
}

// Минимальный вывод
func Print(results []ModuleResult) {
	for _, mr := range results {
		fmt.Println("MODULE:", mr.ModulePath)
		fmt.Println("FILE:  ", mr.ModFilePath)

		if mr.Err != nil {
			fmt.Println("ERROR:", mr.Err)
			fmt.Println()
			continue
		}

		for _, st := range mr.Statuses {
			if st.Error != nil {
				fmt.Printf("  %s  ERROR: %v\n", st.Path, st.Error)
				continue
			}

			if st.NeedUpdate {
				fmt.Printf("  %s  %s → %s (%s)\n",
					st.Path, st.Current, st.Latest, st.UpdateType)
			} else {
				fmt.Printf("  %s  %s (ok)\n", st.Path, st.Current)
			}
		}

		fmt.Println()
	}
}
