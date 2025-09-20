package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"xdvfs/parser"
)

// === Интерактивный режим запуска REPL ===
func (sh *Shell) RunREPL() {
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
		tokens, err := parser.ParseArgs(line)
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
