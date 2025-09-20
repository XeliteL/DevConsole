package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"xdvfs/parser"
)

// === Запуск REPL по скрипту ===
func (sh *Shell) RunScript(path string) error {
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
		tokens, err := parser.ParseArgs(line)
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
