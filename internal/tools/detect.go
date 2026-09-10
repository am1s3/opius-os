package tools

import (
	"os/exec"
	"runtime"
)

type Tool struct {
	Name      string
	Path      string
	Version   string
	Installed bool
}

func DetectTool(name string) Tool {
	path, err := exec.LookPath(name)
	if err != nil {
		return Tool{Name: name, Installed: false}
	}
	
	// Try to get version
	var version string
	cmd := exec.Command(path, "--version")
	if output, err := cmd.Output(); err == nil {
		version = string(output)
		// Truncate to first line
		for i, c := range version {
			if c == '\n' {
				version = version[:i]
				break
			}
		}
	}
	
	return Tool{
		Name:      name,
		Path:      path,
		Version:   version,
		Installed: true,
	}
}

func DetectAll() []Tool {
	tools := []string{
		"espflash",
		"esptool",
		"avrdude",
		"openocd",
		"dfu-util",
		"picotool",
		"git",
		"curl",
	}
	
	result := make([]Tool, 0, len(tools))
	for _, name := range tools {
		result = append(result, DetectTool(name))
	}
	return result
}

func GetInstallCommand(toolName, backend string) string {
	if runtime.GOOS != "darwin" {
		return ""
	}
	
	switch backend {
	case "macports":
		switch toolName {
		case "avrdude":
			return "sudo port install avrdude"
		case "openocd":
			return "sudo port install openocd"
		case "dfu-util":
			return "sudo port install dfu-util"
		case "git":
			return "sudo port install git"
		case "curl":
			return "sudo port install curl"
		}
	case "homebrew":
		switch toolName {
		case "avrdude":
			return "brew install avrdude"
		case "openocd":
			return "brew install openocd"
		case "dfu-util":
			return "brew install dfu-util"
		case "picotool":
			return "brew install picotool"
		case "git":
			return "brew install git"
		case "curl":
			return "brew install curl"
		}
	}
	
	return ""
}