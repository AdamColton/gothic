package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
)

var SerialHelperPackage = "serialHelpers"
var shp *gothicgo.Package

func serialHelperPackage() *gothicgo.Package {
	if shp == nil {
		shp = gothicgo.NewPackage(SerialHelperPackage)
		gothicgo.AutoResolver().Add(SerialHelperPackage, shp.ImportPath)
	}
	return shp
}
