// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Palm1r/QtCToys/qtcreator"
)

func main() {
	infoCommand := flag.NewFlagSet("info", flag.ExitOnError)

	showQtCreator := infoCommand.Bool("qtcreator", false, "Show Qt Creator information")
	showPlugins := infoCommand.Bool("plugins", false, "Show all Qt Creator plugins")
	pluginFilter := infoCommand.String("plugin", "", "Show information for a specific plugin")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'info' command")
		fmt.Println("Usage: qtctoys info [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "info":
		infoCommand.Parse(os.Args[2:])
	case "version":
		fmt.Println("QtCToys v0.1.0")
		return
	case "help":
		fmt.Println("Usage: qtctoys [command] [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  info      Show information about Qt Creator and plugins")
		fmt.Println("  version   Show version")
		fmt.Println("  help      Show help")
		fmt.Println("\nFor command-specific help, use: qtctoys [command] --help")
		return
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Run 'qtctoys help' for usage")
		os.Exit(1)
	}

	if infoCommand.Parsed() {
		info, err := qtcreator.GetInfo()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if !*showQtCreator && !*showPlugins && *pluginFilter == "" {
			*showQtCreator = true
		}

		if *showQtCreator {
			fmt.Printf("Qt Creator Version: %s\n", info.Version)
			fmt.Printf("Qt Creator Path: %s\n", info.Path)
		}

		if *showPlugins || *pluginFilter != "" {
			if *showQtCreator {
				fmt.Println() // Добавляем пустую строку для лучшей читаемости
			}

			fmt.Println("Plugins:")
			handlePlugins(info, *pluginFilter, *showPlugins)
		}
	}
}

func handlePlugins(info qtcreator.Info, pluginName string, showAllPlugins bool) {
	var pluginNames []string
	for name := range info.Plugins {
		pluginNames = append(pluginNames, name)
	}
	sort.Strings(pluginNames)

	if showAllPlugins || pluginName == "all" {
		for _, name := range pluginNames {
			plugin := info.Plugins[name]
			fmt.Printf("  %s %s: %s\n", name, plugin.Version, plugin.Description)
		}
		return
	}

	if pluginName != "" {
		plugin, exists := info.Plugins[pluginName]
		if !exists {
			found := false
			for _, name := range pluginNames {
				if strings.EqualFold(name, pluginName) {
					plugin = info.Plugins[name]
					pluginName = name
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("  Plugin '%s' not found\n", pluginName)

				fmt.Println("  Available plugins containing this term:")
				matchFound := false
				for _, name := range pluginNames {
					if strings.Contains(strings.ToLower(name), strings.ToLower(pluginName)) {
						fmt.Printf("    %s\n", name)
						matchFound = true
					}
				}

				if !matchFound {
					fmt.Println("    No matching plugins found")
				}

				os.Exit(1)
			}
		}

		fmt.Printf("  %s:\n", pluginName)
		fmt.Printf("    Version: %s\n", plugin.Version)
		fmt.Printf("    Description: %s\n", plugin.Description)
	}
}
