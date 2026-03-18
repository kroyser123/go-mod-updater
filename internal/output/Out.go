package output

import (
	"fmt"

	"github.com/kroyser123/go-mod-updater/internal/version"
)

// Красивый CLI-вывод

type ModuleResult struct {
	Path      string
	Module    string
	GoVersion string
	Statuses  []version.Dependecies
	Err       error
}

func NewModuleResult(path, module string, goversion string, statuses []version.Dependecies, err error) ModuleResult {
	return ModuleResult{
		Path:      path,
		Module:    module,
		GoVersion: goversion,
		Statuses:  statuses,
		Err:       err,
	}
}

func Print(results []ModuleResult) {
	for _, r := range results {
		fmt.Printf("\nMODULE: %s\n", r.Module)
		fmt.Printf("FILE:   %s\n", r.Path)
		fmt.Printf("GO:     %s\n", r.GoVersion)

		if r.Err != nil {
			fmt.Printf("ERROR: %v\n", r.Err)
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
				fmt.Printf("%s %s \n", name, st.Current)
			}
		}
		fmt.Println("--------------------------------------------------------------")
	}
}
