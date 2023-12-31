package mplaces

import (
	"GoShopper/internal/model"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

func ProcessURL(url string, db *gorm.DB) {
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

	// После декодирования JSON, сохраняем информацию о товаре в базе данных
	if err := saveProductToDB(&product, db); err != nil {
		fmt.Println("Ошибка при сохранении товара в базе данных:", err)
	}

	//// Выводим данные из структуры
	fmt.Println("ID:", product.ID)
	fmt.Println("Name:", product.Name)
	fmt.Println("SalePrices: ", product.SalePriceU)

	// Раз в 12 часов обновляем информацию о товаре в базе данных
	ticker := time.NewTicker(20 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := updateProductInfo(&product, db); err != nil {
					fmt.Println("Ошибка при обновлении информации о товаре в базе данных:", err)
				}
			}
		}
	}()

	// Ожидаем завершения обработки и обновления для данного URL
	time.Sleep(10 * time.Minute)

	// Завершаем горутину с обновлением информации
	done <- true
}
func saveProductToDB(product *model.Product, db *gorm.DB) error {
	// Проверяем, есть ли уже такой товар в базе данных
	var existingProduct model.Product
	result := db.Where("id = ?", product.ID).First(&existingProduct)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// Если товар уже есть в базе данных, обновляем его информацию
	if result.RowsAffected > 0 {
		existingProduct.Name = product.Name
		existingProduct.SalePriceU = product.SalePriceU
		return db.Save(&existingProduct).Error
	}

	// Если товара нет в базе данных, создаем новую запись
	return db.Create(product).Error
}

// Функция для обновления информации о товаре
func updateProductInfo(product *model.Product, db *gorm.DB) error {
	// Выполняем GET запрос для обновления информации о товаре
	response, err := http.Get(fmt.Sprintf("https://card.wb.ru/cards/detail?appType=1&curr=rub&dest=-1257786&regions=80,38,83,4,64,33,68,70,30,40,86,75,69,22,1,31,66,110,48,71,114&spp=0&nm=%d", product.ID))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Читаем содержимое ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Создаем переменную для декодирования JSON
	var data map[string]interface{}

	// Декодируем JSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	// Проверяем наличие поля "data" в JSON
	if _, ok := data["data"]; !ok {
		return fmt.Errorf("поле 'data' отсутствует в JSON")
	}

	// Получаем нужные поля из JSON и сохраняем в структуру Product
	productData, ok := data["data"].(map[string]interface{})["products"].([]interface{})
	if !ok || len(productData) == 0 {
		return fmt.Errorf("неверный формат JSON или отсутствуют данные о товаре")
	}

	productInfo, ok := productData[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("неверный формат JSON или отсутствуют данные о товаре")
	}

	product.Name = productInfo["name"].(string)
	product.SalePriceU = int64(productInfo["salePriceU"].(float64))

	// Сохраняем обновленную информацию в базе данных
	return saveProductToDB(product, db)
}
