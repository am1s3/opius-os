package components

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/opius-os/opius/internal/ui/style"
)

func PromptString(label, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", style.Info.Render(label), style.Dim.Render(defaultVal))
	} else {
		fmt.Printf("%s: ", style.Info.Render(label))
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func PromptChoice(label string, options []string) string {
	fmt.Println()
	PrintSection(label)
	for i, opt := range options {
		fmt.Printf("  [%d] %s\n", i+1, style.Value.Render(opt))
	}
	fmt.Println()
	for {
		input := PromptString("Type number or name", "")
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			continue
		}
		// Check if it's a number
		if len(input) == 1 && input[0] >= '1' && input[0] <= '0'+byte(len(options)) {
			idx := int(input[0] - '1')
			return options[idx]
		}
		// Check if it's a name
		for _, opt := range options {
			if strings.ToLower(opt) == input {
				return opt
			}
		}
		fmt.Println(BadgeWarn(), "Invalid choice, try again")
	}
}

func PromptConfirm(label string, defaultYes bool) bool {
	defStr := "y/N"
	if defaultYes {
		defStr = "Y/n"
	}
	for {
		input := PromptString(fmt.Sprintf("%s [%s]", label, defStr), "")
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			return defaultYes
		}
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println(BadgeWarn(), "Please type 'y' or 'n'")
	}
}