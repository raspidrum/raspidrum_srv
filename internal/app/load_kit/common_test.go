package loadkit

import (
	"path"
	"runtime"
)

func getDBPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../db/")
	return dir
}

func getProjectPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../")
	return dir
}
