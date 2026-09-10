package components

import (
	"fmt"

	"github.com/opius-os/opius/internal/ui/style"
)

func PrintSection(title string) {
	fmt.Println()
	fmt.Println(style.Section.Render("  "+title))
}