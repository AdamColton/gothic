package stdlib

import (
	"github.com/adamcolton/gothic/gothicgo"
)

var timePkg = gothicgo.MustPackageRef("time")
var Time = struct {
	Pkg  gothicgo.PackageRef
	Time gothicgo.StructType
}{
	Pkg:  timePkg,
	Time: gothicgo.DefStruct(timePkg, "Time"),
}
