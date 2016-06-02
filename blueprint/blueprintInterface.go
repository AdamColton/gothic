package blueprint

import (
	"os"
)

type Blueprint interface {
	Prepare()
	Export()
	Package() string
}

type Generator interface {
	Prepare()
	Export() string
}

var blueprints = map[string][]Blueprint{}

func Register(bp Blueprint) {
	bps, ok := blueprints[bp.Package()]
	if !ok {
		bps = []Blueprint{}
		blueprints[bp.Package()] = bps
	}
	blueprints[bp.Package()] = append(bps, bp)
}

func Export() {
	os.Mkdir(AppPath, 0777)

	for _, bps := range blueprints {
		for _, bp := range bps {
			bp.Prepare()
		}
	}

	for pkg, bps := range blueprints {
		os.Mkdir(AppPath+pkg, 0777)
		for _, bp := range bps {
			bp.Export()
		}
	}
}
