Так у проекта должен быть сборщик make который устанговит драйвер браузера 
так же не забыть про импорт библиотек 
регулярное выражение я подготовил 

-wildberries /https:\/\/card\.wb\.ru\/cards\/detail\?appType=1&curr=rub&dest=-?\d+&regions=(?:\d+(?:,\d+)*,?)*&spp=0&nm=\d+(?:;\d+)*
/gm

-... .

"id"
"name"
"salePriceU"

товар на странице может быть не один а массив похожих товаров разных по цене надо справшивать у юзера в чате что из этого отслеживать
internal/config/cfg.go
internal/db/db.go
internal/model/struct.go
internal/mplaces/wb.go
internal/repositories/gormdb/repo.go
internal/repositories/interface.go
