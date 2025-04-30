package preset

import (
	"fmt"
	"path"
	"runtime"

	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
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

func connectSampler() (*lsampler.LinuxSampler, error) {
	lsClient := lscp.NewClient("localhost", "8888", "1s")

	err := lsClient.Connect()
	ls := lsampler.LinuxSampler{
		Client: lsClient,
		Engine: "sfz",
	}
	if err != nil {
		return &ls, fmt.Errorf("Failed connect to LinuxSampler: %v", err)
	}
	return &ls, nil
}
