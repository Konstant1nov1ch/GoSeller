package main

import (
	"GoShopper/internal/config"
	bootstrap "GoShopper/internal/db"
	"GoShopper/internal/mplaces"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
)

func main() {
	// Вводим URL
	url := ""
	fmt.Scan(&url)
	url = strings.Trim(url, " ")

	// Проверяем наличие "www.wildberries.ru/catalog" во введенной строке
	if strings.Contains(url, "www.wildberries.ru/catalog") {
		// Проверяем доступность URL с помощью HEAD-запроса
		if !checkURLAvailability(url) {
			fmt.Println("Ошибка: URL недоступен.")
			return
		}

		// Создаем объект ConfigProvider
		cp, err := config.NewConfigProvider()
		if err != nil {
			fmt.Println("Ошибка при инициализации настроек:", err)
			return
		}

		// Инициализируем конфигурацию, включая Playwright
		config.Init()

		// Подключаемся к базе данных
		db, err := bootstrap.InitGormDB(cp.GetConfig())
		if err != nil {
			fmt.Println("Ошибка подключения к базе данных:", err)
			return
		}

		// Запускаем обработку URL
		processMultipleURLs(url, db)
	} else {
		fmt.Println("Ошибка: Введенный URL не соответствует требуемому формату.")
	}
}

func processMultipleURLs(url string, db *gorm.DB) {
	// Разбиваем URL по запятым, чтобы обрабатывать несколько товаров
	urls := strings.Split(url, ",")
	var wg sync.WaitGroup

	for _, u := range urls {
		u = strings.TrimSpace(u)
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			mplaces.ProcessURL(u, db)
		}(u)
	}

	wg.Wait()
}

// Функция для проверки доступности URL с помощью HEAD-запроса
func checkURLAvailability(url string) bool {
	// Добавляем протокол, если его нет
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Выполняем HEAD-запрос
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Проверяем статусный код ответа
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
