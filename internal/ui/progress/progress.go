package progress

import (
	"fmt"
	"strings"
)

type Bar struct {
	total     int64
	current   int64
	label     string
	width     int
	lastPct   int
}

func NewBar(total int64, label string) *Bar {
	return &Bar{
		total:   total,
		current: 0,
		label:   label,
		width:   40,
		lastPct: -1,
	}
}

func (b *Bar) Update(current int64) {
	b.current = current

	if b.total <= 0 {
		fmt.Printf("\r  %s: %s", b.label, formatBytesSimple(current))
		return
	}

	pct := int(float64(current) / float64(b.total) * 100)
	if pct == b.lastPct {
		return
	}
	b.lastPct = pct

	filled := int(float64(b.width) * float64(pct) / 100)
	if filled > b.width {
		filled = b.width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", b.width-filled)

	fmt.Printf("\r  %s [%s] %3d%% %s",
		b.label,
		bar,
		pct,
		formatBytesSimple(current)+"/"+formatBytesSimple(b.total),
	)
}

func (b *Bar) Finish() {
	if b.total > 0 {
		b.Update(b.total)
	}
	fmt.Println()
}

func formatBytesSimple(bytes int64) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}