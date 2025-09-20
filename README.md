1. Общее описание
   Проект представляет собой эмулятор командной оболочки UNIX-подобной ОС, написанный на языке Go.

   Эмулятор работает в двух режимах:
      1) REPL (интерактивный ввод команд).
      2) Скриптовый режим (чтение и выполнение команд из файла).
   
   Основные возможности:
      1) Построение виртуальной файловой системы (VFS) в памяти.
      2) Поддержка базовых команд (ls, cd, exit) и дополнительных (history, uname, rmdir).
      3) Обработка ошибок (неизвестные команды, неверные аргументы, операции с файлами/директориями).
   
   Проект выполнен в 5 этапов:
      1) REPL с заглушками.
      2) Поддержка параметров запуска (--VFS, --Script).
      3) Загрузка и хранение VFS в памяти.
      4) Реализация команд ls, cd, history, uname.
      5) Реализация команды rmdir (удаление пустой директории).

2. Установка и запуск 
   Требуется GO(1.20+)

   Клонируем проект из гита:
      git clone <repo_url>
      cd project

   Запуск в интерактивном режиме: 
      go run main.go --VFS ./data
   Запуск со скриптом:
      go run main.go --VFS ./data --Script ./scripts/stage4.txt

3. Список команд:
	   ls [path]  - показ содержимого
	   cd <path>  - смена пути
	   history    - показ истории команд
	   uname      - показ информация об ОС
	   rmdir <d>  - удаление пустой директории
	   exit, quit - выход из эмулятора

4. Тестовые скрипты, находящиеся в папке scripts:
   # Этап 1
   go run main.go --Script ./scripts/stage1.txt

   # Этап 2
   go run main.go --VFS ./data --Script ./scripts/stage2.txt

   # Этап 3
   go run main.go --VFS ./data --Script ./scripts/stage3_nested.txt

   # Этап 4
   go run main.go --VFS ./data --Script ./scripts/stage4.txt

   # Этап 5
   go run main.go --VFS ./data --Script ./scripts/stage5.txt

5. Пример структуры vfs:
      data/
         file1.txt
         file2.txt
         subdir/
            nested.txt
         subdir_empty/