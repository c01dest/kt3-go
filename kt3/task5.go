package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"
)

const numWorkers = 4

var outMu sync.Mutex

func processImage(path string) {
	// 1. Проверяем, существует ли файл и обычный ли он (не папка, не симлинк)
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return
	}

	// 2. Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// 3. Считываем только заголовок (разрешение) без декодирования пикселей
	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	// 4. Формируем тег и атомарно выводим целой строкой
	tag := fmt.Sprintf(`<img src="%s" width="%d" height="%d" />`, filepath.Base(path), cfg.Width, cfg.Height)

	outMu.Lock()
	fmt.Println(tag)
	outMu.Unlock()
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		return
	}

	jobs := make(chan string)
	var wg sync.WaitGroup

	// Запускаем фиксированный пул воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				processImage(path)
			}
		}()
	}

	// Передаем пути к файлам в канал задач
	for _, path := range files {
		jobs <- path
	}
	close(jobs)

	wg.Wait()
}