package parser

import (
	"errors"
	"strings"
)

// Вспомогательная функция для обработки одного символа
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
