package components

import (
	"github.com/opius-os/opius/internal/ui/style"
)

func BadgeOK() string   { return style.Success.Render("[ OK ]") }
func BadgeWarn() string { return style.Warning.Render("[WARN]") }
func BadgeErr() string  { return style.Error.Render("[FAIL]") }
func BadgeInfo() string { return style.Info.Render("[INFO]") }