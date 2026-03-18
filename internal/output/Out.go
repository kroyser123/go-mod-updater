package output

import (
	"fmt"

	"github.com/kroyser123/go-mod-updater/internal/version"
)

// Красивый CLI-вывод

type ModuleResult struct {
	Path     string
	Module   string
	Statuses []version.Dependecies
	Err      error
}

func NewModuleResult(path, module string, statuses []version.Dependecies, err error) ModuleResult {
	return ModuleResult{
		Path:     path,
		Module:   module,
		Statuses: statuses,
		Err:      err,
	}
}

func Print(results []ModuleResult) {
	for _, r := range results {
		fmt.Printf("\nMODULE: %s\n", r.Module)
		fmt.Printf("FILE:   %s\n\n", r.Path)

		if r.Err != nil {
			fmt.Printf("ERROR: %v\n\n", r.Err)
			continue
		}
		for _, st := range r.Statuses {
			if st.Path == "" {
				continue
			}
			name := st.Path
			if st.Indirect {
				name += " (indirect)"
			}
			if st.Error != nil {
				fmt.Printf("  %s ERROR: %v\n", name, st.Error)
				continue
			}
			if st.NeedUpdate {
				fmt.Printf("  %s %s → %s (%s)\n",
					name, st.Current, st.Latest, st.UpdateType)
			} else {
				fmt.Printf("%s %s %v \n", name, st.Current, st.NeedUpdate)
			}
		}
		fmt.Println("----------------------------------------------")
	}
}
