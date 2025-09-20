package main

import (
	"flag"
	"fmt"

	"xdvfs/shell"
	"xdvfs/vfs"
)

func main() {
	// === Обработка параметров командной строки ===
	pathVFS := flag.String("VFS", "", "Путь к директории VFS")
	pathScript := flag.String("Script", "", "Путь к стартовому скрипту")
	flag.Parse()

	// Отладочный вывод параметров
	fmt.Println("=== Запуск эмулятора ===")
	fmt.Printf("VFS path: %s\n", *pathVFS)
	fmt.Printf("Script path: %s\n", *pathScript)

	// Создание VFS
	var root *vfs.Node
	var err error
	if *pathVFS != "" {
		root, err = vfs.BuildVFS(*pathVFS, nil)
		if err != nil {
			fmt.Printf("Ошибка при загрузке VFS: %v\n", err)
			root = vfs.NewEmptyVFS()
		}
	} else {
		root = vfs.NewEmptyVFS()
	}

	sh := shell.NewShell(root) // Создаём оболочку с корнем VFS

	// Запуск скрипта
	if *pathScript != "" {
		if err := sh.RunScript(*pathScript); err != nil {
			fmt.Printf("Скрипт остановлен: %v\n", err)
			return
		}
	}

	// Запуск в интерактивном режиме при отсутствии скрипта
	sh.RunREPL()
}
