package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Структура VFS
type VFS struct {
	Name string
}

// === Создание нового VFS ===
func newVFS(name string) *VFS {
	if name == "" {
		name = "userVFS" // Имя по умолчанию
	}
	return &VFS{Name: name}
}

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
			// Пропуск множественных пробелов
		}
	default:
		curToken.WriteRune(r)
	}

	// Вывод ошибки, если в конце строки служебный символ
	if isLast && escaped {
		return nil, inSingle, inDouble, escaped, errors.New("Служебный символ в конце строки")
	}

	return args, inSingle, inDouble, escaped, nil
}

// === Парсер строки, вводимой пользователем, с экранированием и поддержкой кавычек ===
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
		return nil, errors.New("Незакрытая кавычка")
	}

	// Добавление последнего токена при его наличии
	if curToken.Len() > 0 {
		args = append(args, curToken.String())
	}

	return args, nil
}

// === Обработчик команд ===
func handleCommand(vfs *VFS, args []string) {
	// При пустой строке
	if len(args) == 0 {
		return
	}

	cmd := args[0] // Определение команды
	switch cmd {
	case "help":
		printHelp()
	case "switch", "quit":
		fmt.Println("Выход")
		os.Exit(0)
	case "ls":
		// ls может выводить максимум одну директорию
		if len(args) > 2 {
			fmt.Println("Ошибка: неправильное использование ls [path]")
			return
		}
		fmt.Printf("[stub] ls вызван. Аргументы: %v/n", args[1:]) // Заглушка
	case "cd":
		// У cd может быть только 1 аргумент
		if len(args) != 2 {
			fmt.Println("Ошибка: неправильное использование cd <path>")
			return
		}
		fmt.Printf("[stub] cd вызван. Аргумент: %s/n", args[1]) // Заглушка
	default:
		fmt.Printf("Ошбика: неизвестная команда '%s'\n", cmd)
	}
}

// === Список команд ===
func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println(" help       - показать подсказку")
	fmt.Println(" ls [path]  - показать содержимое(заглушка)")
	fmt.Println(" cd <path>  - сменить путь(заглушка)")
	fmt.Println(" exit, quit - выйти из эмулятора")
}

func main() {
	// === Инициализация эмулятора ===
	vfs := newVFS("xdVFS") // Создание VFS

	// Приветственное сообщение
	fmt.Printf("Эмулятор оболочки. VFS: %s\n", vfs.Name)
	fmt.Println("Введите 'help' для списка команд. Для выхода введите команды: 'quit' либо 'exit'.")

	scanner := bufio.NewScanner(os.Stdin) // Сканер для чтения строк

	// === Главный цикл REPL ===
	for {
		fmt.Printf("%s> ", vfs.Name) // Приглашение к вводу

		// Ошибка ввода или EOF
		if !scanner.Scan() {
			fmt.Println()
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue // Пропуск пустого ввода
		}

		tokens, err := parseArgs(line) // Разбиение строки на токены
		if err != nil {                // Ошибка парсинга
			fmt.Printf("Ошибка парсера: %v\n", err)
			continue
		}

		if len(tokens) == 0 {
			continue
		}

		handleCommand(vfs, tokens) // Передача токенов обработчику команд
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения ввода %v\n", err)
	}
}
