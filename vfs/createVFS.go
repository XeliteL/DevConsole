package vfs

import (
	"os"
	"path/filepath"
)

// Создание нового пустого VFS
func NewEmptyVFS() *Node {
	return &Node{
		Name:  "xdVFS",
		IsDir: true,
	}
}

// Построение дерева из директории
func BuildVFS(path string, parent *Node) (*Node, error) {
	info, err := os.Stat(path) // Получение информации о файле/директории
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name:   info.Name(),
		IsDir:  info.IsDir(),
		Parent: parent,
	}

	// Чтение содержимого, если это директория
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		// Создание дочернего узла для каждого элемента
		for _, e := range entries {
			childPath := filepath.Join(path, e.Name())
			child, err := BuildVFS(childPath, node)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		}
	}

	return node, nil
}
