package flash

import "context"

type BinaryEntry struct {
	Path    string `json:"path"`
	Address string `json:"address"`
}

type Target struct {
	Chip     string        `json:"chip"`
	Board    string        `json:"board"`
	Port     string        `json:"port"`
	File     string        `json:"file"`
	Binaries []BinaryEntry `json:"binaries"`
	Baud     int           `json:"baud"`
	Erase    bool          `json:"erase"`
}

type Step struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Command     []string `json:"command"`
	Phase       string   `json:"phase"`
}

type Plan struct {
	Target  Target `json:"target"`
	Driver  string `json:"driver"`
	Tool    string `json:"tool"`
	Steps   []Step `json:"steps"`
	Summary string `json:"summary"`
}

type ProgressEvent struct {
	Step    string `json:"step"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
	Percent int    `json:"percent"`
}

type Driver interface {
	Name() string
	Supports(t Target) bool
	Plan(ctx context.Context, t Target) (*Plan, error)
	Execute(ctx context.Context, p *Plan, dryRun bool, ev chan<- ProgressEvent) error
}