package main

import (
	"GoShopper/internal/config"
	"GoShopper/internal/mplaces"
	"fmt"
	"net/http"
	"strings"
)

func main() {
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

		// Инициализируем конфигурацию, включая Playwright
		config.Init()

		// Запускаем обработку URL
		mplaces.ProcessURL(url)
	} else {
		fmt.Println("Ошибка: Введенный URL не соответствует требуемому формату.")
	}
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