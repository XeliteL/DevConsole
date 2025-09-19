package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// === Структуры данных для VFS ===

// Папка/файл внутри VFS
type Node struct {
	Name     string // Имя файла/папки
	IsDir    bool // Файл или папка
	Children []*Node // Список дочерних элементов
	Parent   *Node // Ссылка на родительскую директорию
}

// Текущая директория
var currentDir *Node

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
func handleCommand(args []string) error {
	// При пустой строке
	if len(args) == 0 {
		return nil
	}

	cmd := args[0] // Определение команды
	switch cmd {
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Выход")
		os.Exit(0)
	case "ls":
		// ls может выводить максимум одну директорию
		if len(args) > 2 {
			return errors.New("ошибка: неправильное использование ls [path]")
		}
		fmt.Printf("[stub] ls вызван. Аргументы: %v\n", args[1:]) // Заглушка
	case "cd":
		// У cd может быть только 1 аргумент
		if len(args) != 2 {
			return errors.New("ошибка: неправильное использование cd <path>")
		}
		fmt.Printf("[stub] cd вызван. Аргумент: %s\n", args[1]) // Заглушка
	default:
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
	fmt.Println("  exit, quit - выйти из эмулятора")
}

// === Интерактивный режим запуска REPL ===
func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", currentDir.Name) // Приглашение к вводу

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
		if err := handleCommand(tokens); err != nil {
			fmt.Printf("Ошибка обработчика: %v\n", err)
		}
	}

	// Проверка ошибок ввода
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
	}
}

// === Запуск REPL по скрипту ===
func runScript(path string) error {
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
		fmt.Printf("%s> %s\n", currentDir.Name, line)

		// Парсинг и выполнение команды обработчика
		tokens, err := parseArgs(line)
		if err != nil {
			fmt.Printf("Ошибка парсера: %v\n", err)
			continue
		}

		if err := handleCommand(tokens); err != nil {
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
	currentDir = root // Ставаим корень, как текущую директорию

	// Запуск скрипта
	if *pathScript != "" {
		if err := runScript(*pathScript); err != nil {
			fmt.Printf("Скрипт остановлен: %v\n", err)
			return
		}
	}

	// Запуск в интерактивном режиме при отсутствии скрипта
	runREPL()
}