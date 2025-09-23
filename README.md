1. Проект
   Проект представляет собой эмулятор командной оболочки UNIX-подобной ОС, написанный на языке Go.

   Эмулятор поддерживает два режима работы:
      1) Интерактивный режим (REPL) — пользователь вводит команды вручную.
      2) Скриптовый режим — выполнение команд из файла (имитация диалога с пользователем).

   Виртуальная файловая система (VFS) хранится в памяти и может загружаться из реальной директории.

   Проект реализован поэтапно (от простого REPL до полнофункционального эмулятора с командами ls, cd, history, uname, rmdir).

   Структура проекта: 
      devConsole
      ├── bats/                   # Запуск скриптов
      ├── data/                   # Файлы данных и конфигурации
      ├── scripts/                # Тестовые скрипты
      ├── parser/                 # Модуль анализа токена
      │   ├── parser.go
      │   └── switchScreening.go
      ├── shell/                  # Модуль командной оболочки
      │   ├── commands.go
      │   ├── help.go
      │   ├── interactiveREPL.go
      │   ├── scriptREPL.go
      │   └── shell.go
      ├── vfs/                    # Модуль виртуальной файловой системы
      │   ├── createVFS.go
      │   └── vfs.go
      └── main.go                 # Главная программа

2. Функции и настройки
   Основные функции:
      1) buildVFS(path string, parent *Node) — загружает дерево директорий с диска.
      2) newEmptyVFS() — создаёт пустую VFS в памяти.
      3) parseArgs(line string) — парсинг аргументов команд с кавычками и экранированием.
      4) handleCommand(args []string) — обработка команд (ls, cd, history, uname, rmdir).
      5) runREPL() — интерактивный режим.
      6) runScript(path string) — запуск команд из файла-скрипта.

   Параметры запуска:
      1) --VFS <path> — путь к директории, из которой строится VFS.
      2) --Script <path> — путь к файлу со скриптом для выполнения.

   Поддерживаемые команды:
      1) help        - вывести список доступных команд
      2) ls [path]   - показать содержимое директории
      3) cd <dir>    - перейти в директорию (cd .. — на уровень выше)
      4) history     - вывести историю команд
      5) uname       - показать информацию о ОС и архитектуре
      6) rmdir <dir> - удалить пустую директорию
      7) exit, quit  - завершить работу

3. Сборка проекта и запуск тестов
   Сборка проекта:
      go build -o emulator main.go

   Запуск вручную
   # Запуск в REPL с тестовой VFS:
      go run main.go --VFS ./data

   # Запуск со скриптом:
      go run main.go --VFS ./data --Script ./scripts/stage4_commands.txt

   Запуск готовых тестов осуществляется в папке bats:
      bats/
      ├── stage1.bat       # тест Этапа 1
      ├── stage2.bat       # тест Этапа 2
      ├── stage3.bat       # тест Этапа 3
      ├── stage4.bat       # тест Этапа 4
      └── stage5.bat       # тест Этапа 5

4. Примеры использования:
   Пример 1. Навигация по директориям:
      data> ls
      [DIR] subdir
            file1.txt
            file2.txt

      data> cd subdir
      subdir> ls
            nested.txt

   Пример 2. История команд и uname:
      data> ls
      [DIR] subdir
            file1.txt

      data> uname
      OS: windows
      Arch: amd64

      data> history
      1: ls
      2: uname
      3: history

   Пример 3. Удаление пустой директории:
      data> ls
      [DIR] subdir_empty
            file1.txt

      data> rmdir subdir_empty
      Каталог 'subdir_empty' удалён

      data> ls
            file1.txt