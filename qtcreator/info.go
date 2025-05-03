package qtcreator

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Info struct {
	Version string // Qt Creator version
	Path    string // Path to the executable
}

func GetInfo() (Info, error) {
	var info Info
	var creatorPath string
	var err error

	switch runtime.GOOS {
	case "windows":
		possiblePaths := []string{
			"C:\\Qt\\Tools\\QtCreator\\bin\\qtcreator.exe",
			"C:\\Program Files\\Qt\\Tools\\QtCreator\\bin\\qtcreator.exe",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				creatorPath = path
				break
			}
		}

	case "darwin":
		homeDir, _ := os.UserHomeDir()
		creatorPath = homeDir + "/Qt/Qt Creator.app/Contents/MacOS/Qt Creator"

	default:
		creatorPath, err = exec.LookPath("qtcreator")
	}

	if creatorPath == "" {
		return info, fmt.Errorf("could not find Qt Creator")
	}

	info.Path = creatorPath

	cmd := exec.Command(creatorPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return info, fmt.Errorf("error when requesting version: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Qt Creator") {
			parts := strings.SplitN(line, "Qt Creator", 2)
			if len(parts) > 1 {
				info.Version = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if info.Version == "" {
		return info, fmt.Errorf("could not determine Qt Creator version")
	}

	return info, nil
}
