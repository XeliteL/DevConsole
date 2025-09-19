Эмулятор консоли разработчика для конфигурационного управления

Команды для запусков скриптов: 
1) go run main.go --VFS ./data --Script ./scripts/stage1.txt   
2) go run main.go --VFS ./data --Script ./scripts/stage2.txt
3) 1) go run main.go --VFS ./data --Script ./scripts/stage3_empty.txt      
   2) go run main.go --VFS ./data --Script ./scripts/stage3_files.txt
   3) go run main.go --VFS ./data --Script ./scripts/stage3_nested.txt      
4) go run main.go --VFS ./data --Script ./scripts/stage4.txt   
5) go run main.go --VFS ./data --Script ./scripts/stage5.txt   