package vfs

// Папка/файл внутри VFS
type Node struct {
	Name     string  // Имя файла/папки
	IsDir    bool    // Файл или папка
	Children []*Node // Список дочерних элементов
	Parent   *Node   // Ссылка на родительскую директорию
}
