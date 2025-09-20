package shell

import (
	"fmt"
)

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
