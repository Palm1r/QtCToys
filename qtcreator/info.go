package qtcreator

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type PluginInfo struct {
	Name        string // Plugin name
	Version     string // Plugin version
	Description string // Plugin description
}

type QtCreatorInfo struct {
	Version string                // Qt Creator version
	Path    string                // Path to the executable
	Plugins map[string]PluginInfo // Available plugins
}

func GetInfo() (QtCreatorInfo, error) {
	var info QtCreatorInfo
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
	info.Plugins = make(map[string]PluginInfo)

	cmd := exec.Command(creatorPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return info, fmt.Errorf("error when requesting version: %w", err)
	}

	err = parseVersionOutput(string(output), &info)
	if err != nil {
		return info, fmt.Errorf("error parsing version output: %w", err)
	}

	return info, nil
}

func parseVersionOutput(output string, info *QtCreatorInfo) error {
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "Qt Creator") {
			parts := strings.SplitN(line, "Qt Creator", 2)
			if len(parts) > 1 {
				versionParts := strings.Split(strings.TrimSpace(parts[1]), " ")
				info.Version = versionParts[0]
				break
			}
		}

		if i > 5 {
			break
		}
	}

	if info.Version == "" {
		return fmt.Errorf("could not determine Qt Creator version")
	}

	// Process plugins
	var currentPlugin string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		if len(parts) >= 2 {
			if isVersionFormat(parts[1]) {
				currentPlugin = parts[0]
				plugin := PluginInfo{
					Name:    currentPlugin,
					Version: parts[1],
				}

				if len(parts) >= 3 {
					plugin.Description = strings.TrimSpace(parts[2])
				}

				info.Plugins[currentPlugin] = plugin
			}
		}
	}

	return nil
}

func isVersionFormat(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if !isNumeric(part) {
			return false
		}
	}

	return true
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func GetPluginInfo(name string) (PluginInfo, error) {
	info, err := GetInfo()
	if err != nil {
		return PluginInfo{}, err
	}

	plugin, exists := info.Plugins[name]
	if !exists {
		return PluginInfo{}, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin, nil
}
