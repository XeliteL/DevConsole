package parser

import (
	"errors"
	"strings"
)

// Разбиение строки на токены, учитывая кавычки и экранирование
func ParseArgs(line string) ([]string, error) {
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
