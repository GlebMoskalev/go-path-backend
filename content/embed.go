package content

import "embed"

//go:embed all:theory
var TheoryFS embed.FS

//go:embed all:tasks
var TasksFS embed.FS
