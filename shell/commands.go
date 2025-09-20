package shell

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// === Обработчик команд ===
func (sh *Shell) handleCommand(args []string) error {
	// При пустой строке
	if len(args) == 0 {
		return nil
	}

	// Сохранение команды в историю
	sh.history = append(sh.history, strings.Join(args, " "))

	cmd := args[0] // Определение команды
	switch cmd {
	case "help": // Вывод окна помощи
		printHelp()

	case "exit", "quit": // Команды завершения программы
		fmt.Println("Выход")
		os.Exit(0)

	case "ls": // Вывод содержимого директории
		dir := sh.currentDir // При отсутствии аргумента выводим текущую директорию
		// При указании пути
		if len(args) == 2 {
			found := false
			// Поиск среди дочерних элементов директории с нужным именем
			for _, child := range sh.currentDir.Children {
				if child.Name == args[1] && child.IsDir {
					dir = child
					found = true
					break
				}
			}

			// Если папка не найдена
			if !found {
				return fmt.Errorf("каталог '%s' не найден", args[1])
			}
		}

		// Вывод содержимого
		for _, child := range dir.Children {
			if child.IsDir {
				fmt.Printf("[DIR] %s\n", child.Name)
			} else {
				fmt.Printf("      %s\n", child.Name)
			}
		}

	case "cd": // Переход к директории
		// У cd может быть только 1 аргумент
		if len(args) != 2 {
			return errors.New("ошибка: использование cd <path>")
		}

		target := args[1] // Цель

		// Переход к родительской директории
		if target == ".." && sh.currentDir.Parent != nil {
			sh.currentDir = sh.currentDir.Parent
			return nil
		}

		// Поиск дочерней директории с нужным именем
		for _, child := range sh.currentDir.Children {
			if child.Name == target && child.IsDir {
				sh.currentDir = child
				return nil
			}
		}

		// Если директория не найдена
		return fmt.Errorf("ошибка: каталог '%s' не найден", target)

	case "history": // Показ истории вызванных команд
		for i, cmd := range sh.history {
			fmt.Printf("%d: %s\n", i+1, cmd)
		}

	case "uname": // Показ сведений об ОС
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Arch: %s\n", runtime.GOARCH)

	case "rmdir": // Удаление пустой директории
		if len(args) != 2 {
			return errors.New("ошибка: использование rmdir <dirname>")
		}

		dirName := args[1]
		for i, child := range sh.currentDir.Children {
			if child.Name == dirName {
				// Проверка на директорию
				if !child.IsDir {
					return fmt.Errorf("'%s' - это файл, а не директория", dirName)
				}

				// Проверка на заполненность
				if len(child.Children) > 0 {
					return fmt.Errorf("каталог '%s' не пуст", dirName)
				}

				// Удаляем узел
				sh.currentDir.Children = append(sh.currentDir.Children[:i], sh.currentDir.Children[i+1:]...)
				fmt.Printf("Каталог '%s' удалён\n", dirName)
				return nil
			}
		}

		// Если директория не найдена
		return fmt.Errorf("ошибка: каталог '%s' не найден", dirName)

	default: // При неизвестной команде
		return fmt.Errorf("ошибка: неизвестная команда '%s'", cmd)
	}

	return nil
}
