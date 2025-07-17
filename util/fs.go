package util

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func AbsPathify(basePath, inPath string) string {
	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = userHomeDir() + inPath[5:]
	}

	inPath = os.ExpandEnv(inPath)

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	var pt string
	if basePath != "" {
		pt = path.Join(basePath, inPath)
	} else {
		pt = inPath
	}

	p, err := filepath.Abs(pt)
	if err == nil {
		return filepath.Clean(p)
	}

	slog.Error(fmt.Errorf("could not discover absolute path: %w", err).Error())

	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
