package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// === Структуры данных для VFS ===

// Папка/файл внутри VFS
type Node struct {
	Name     string  // Имя файла/папки
	IsDir    bool    // Файл или папка
	Children []*Node // Список дочерних элементов
	Parent   *Node   // Ссылка на родительскую директорию
}

type Shell struct {
	currentDir *Node    // Текущая директория
	history    []string // История команд
}

// Создание нового пустого VFS
func newEmptyVFS() *Node {
	return &Node{
		Name:  "xdVFS",
		IsDir: true,
	}
}

// Построение дерева из директории
func buildVFS(path string, parent *Node) (*Node, error) {
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
			child, err := buildVFS(childPath, node)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		}
	}

	return node, nil
}

// === Парсер строки, вводимой пользователем, с экранированием и поддержкой кавычек ===
func switchScreening(r rune, isLast bool, curToken *strings.Builder,
	args []string, inSingle, inDouble, escaped bool,
) ([]string, bool, bool, bool, error) {
	if escaped {
		// Добавляем символ буквально и сбрасываем флаг
		curToken.WriteRune(r)
		escaped = false
		return args, inSingle, inDouble, escaped, nil
	}

	switch r {
	case '\\': // Включаем экранирование
		escaped = true
	case '\'':
		if inDouble {
			// В режиме двойных кавычек одиночная воспринимаетсябуквлаьно, поэтоьу добавляем её в токен
			curToken.WriteRune(r)
		} else {
			// Выключаем режим одиночной кавычки и не добавляем кавычку в токен
			inSingle = !inSingle
		}
	case '"':
		if inSingle {
			// В режиме одинарных кавычек двойная воспринимается буквлаьно, поэтому добавляем её в токен
			curToken.WriteRune(r)
		} else {
			// Выключаем режим двойной кавычки и не добавляем кавычку в токен
			inDouble = !inDouble
		}
	case ' ', '\t':
		if inSingle || inDouble {
			// Внутри кавычек пробел и таб воспринимаются буквально
			curToken.WriteRune(r)
		} else {
			// Разделитель токенов
			if curToken.Len() > 0 {
				args = append(args, curToken.String())
				curToken.Reset()
			}
		}
	default:
		curToken.WriteRune(r)
	}

	// Вывод ошибки, если в конце строки служебный символ
	if isLast && escaped {
		return nil, inSingle, inDouble, escaped, errors.New("служебный символ в конце строки")
	}

	return args, inSingle, inDouble, escaped, nil
}

// Возвращает слайс токенов либо ошибку при незакрытой кавычки
func parseArgs(line string) ([]string, error) {
	var args []string
	var curToken strings.Builder
	inSingle := false // Одинарные кавычки
	inDouble := false // Двойные кавычки
	escaped := false  // Экранирование

	for i, r := range line {
		var err error
		args, inSingle, inDouble, escaped, err = switchScreening(
			r, i == len(line)-1, &curToken,
			args, inSingle, inDouble, escaped,
		)
		if err != nil {
			return nil, err
		}
	}

	// Проверка закрытия кавычек
	if inSingle || inDouble {
		return nil, errors.New("незакрытая кавычка")
	}

	// Добавление последнего токена при его наличии
	if curToken.Len() > 0 {
		args = append(args, curToken.String())
	}

	return args, nil
}

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
				return fmt.Errorf("каталаог '%s' не найден", args[1])
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

// === Список команд ===
func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  help       - показать подсказку")
	fmt.Println("  ls [path]  - показать содержимое(заглушка)")
	fmt.Println("  cd <path>  - сменить путь(заглушка)")
	fmt.Println("  history    - показать историю команд")
	fmt.Println("  uname      - информация об ОС")
	fmt.Println("  rmdir <d>  - удаление пустой директории")
	fmt.Println("  exit, quit - выйти из эмулятора")
}

// === Интерактивный режим запуска REPL ===
func (sh *Shell) runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", sh.currentDir.Name) // Приглашение к вводу

		// Чтение строки пользоватиля
		if !scanner.Scan() { // Ошибка ввода
			fmt.Println()
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" { // Пропуск пустого ввода
			continue
		}

		// Разбиение строки на токены
		tokens, err := parseArgs(line)
		if err != nil {
			fmt.Printf("Ошибка парсера: %v\n", err)
			continue
		}

		// Выполнение команды обработчика
		if err := sh.handleCommand(tokens); err != nil {
			fmt.Printf("Ошибка обработчика: %v\n", err)
		}
	}

	// Проверка ошибок ввода
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
	}
}

// === Запуск REPL по скрипту ===
func (sh *Shell) runScript(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть скрипт: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Пропуск пустых строк и комментариев
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Имитация диалога с пользователем
		fmt.Printf("%s> %s\n", sh.currentDir.Name, line)

		// Парсинг и выполнение команды обработчика
		tokens, err := parseArgs(line)
		if err != nil {
			fmt.Printf("Ошибка парсера: %v\n", err)
			continue
		}

		if err := sh.handleCommand(tokens); err != nil {
			fmt.Printf("Ошибка обработчика: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения скрипта: %v", err)
	}
	return nil
}

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
	var root *Node
	var err error
	if *pathVFS != "" {
		root, err = buildVFS(*pathVFS, nil)
		if err != nil {
			fmt.Printf("Ошибка при загрузке VFS: %v\n", err)
			root = newEmptyVFS()
		}
	} else {
		root = newEmptyVFS()
	}

	sh := &Shell{currentDir: root} // Ставаим корень, как текущую директорию

	// Запуск скрипта
	if *pathScript != "" {
		if err := sh.runScript(*pathScript); err != nil {
			fmt.Printf("Скрипт остановлен: %v\n", err)
			return
		}
	}

	// Запуск в интерактивном режиме при отсутствии скрипта
	sh.runREPL()
}
