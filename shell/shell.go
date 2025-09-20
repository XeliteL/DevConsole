package shell

import (
	"xdvfs/vfs"
)

// Оболочка с состоянием
type Shell struct {
	currentDir *vfs.Node // Текущая директория
	history    []string  // История команд
}

// Создание оболочки
func NewShell(root *vfs.Node) *Shell {
	return &Shell{currentDir: root}
}
