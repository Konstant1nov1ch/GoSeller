package mplaces

import (
	"GoShopper/internal/model"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	"regexp"
)

func ProcessURL(url string) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Регулярное выражение для поиска искомого URL
	urlPattern := regexp.MustCompile(`https:\/\/card\.wb\.ru\/cards\/detail\?appType=1&curr=rub&dest=-?\d+&regions=(?:\d+(?:,\d+)*,?)*&spp=0&nm=\d+(?:;\d+)*`)

	// Переменная для хранения найденного URL
	foundURL := ""

	// Обработчик события response
	page.On("response", func(response playwright.Response) {
		// Если уже нашли URL, ничего не делаем
		if foundURL != "" {
			return
		}

		// Ищем URL в содержимом ответа
		if matches := urlPattern.FindStringSubmatch(response.URL()); len(matches) > 0 {
			foundURL = matches[0]
			fmt.Printf("Найден URL: %s\n", foundURL)
		}
	})

	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	// Выполняем GET запрос
	response, err := http.Get(foundURL)
	if err != nil {
		fmt.Println("Ошибка при выполнении GET запроса:", err)
		return
	}
	defer response.Body.Close()

	// Читаем содержимое ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении содержимого ответа:", err)
		return
	}

	// Создаем переменную для декодирования JSON
	var data map[string]interface{}

	// Декодируем JSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	// Получаем нужные поля из JSON и сохраняем в структуру Product
	var product model.Product
	productData := data["data"].(map[string]interface{})["products"].([]interface{})[0].(map[string]interface{})
	product.ID = int64(int(productData["id"].(float64)))
	product.Name = productData["name"].(string)
	product.SalePriceU = int64(int(productData["salePriceU"].(float64)))
	//// Выводим данные из структуры
	fmt.Println("ID:", product.ID)
	fmt.Println("Name:", product.Name)
	fmt.Println("SalePrices: ", product.SalePriceU)
}
