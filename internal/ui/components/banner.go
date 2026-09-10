package components

import (
	"fmt"

	"github.com/opius-os/opius/internal/build"
	"github.com/opius-os/opius/internal/config"
	"github.com/opius-os/opius/internal/ui/style"
)

const logo = `
   ___  ___  ___  _  _  _  ___
  / _ \| _ \|_ _|| || | | |/ __|
 | (_) |  _/ | | | __ | | |\__ \
  \___/|_|  |___||_||_||_|_|___/
`

func PrintBanner() {
	fmt.Println(style.Brand.Render(logo))
	fmt.Println(style.Dim.Render("  OPIUS OS · terminal embedded workspace · v" + build.Version + "-alpha"))
	fmt.Println(style.Dim.Render("  " + config.DefaultRepository))
	fmt.Println()
}