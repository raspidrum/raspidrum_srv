package db

import (
	"path"
	"runtime"
)

func getDBPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../db/")
	return dir
}
