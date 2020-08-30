# MMORPG
<img src="https://2ch.hk/gd/src/634039/15814183358363.png" />
<img src="https://habrastorage.org/webt/qe/tt/ll/qettllbkbgqjvfn_o_x4oaosjzw.gif" />
Проект задуман как многопользовательская игра с бесконечно генерируемым миром и изометрической графикой.
Как все это примерно работает описано в моей статье на Хабре:
https://habr.com/ru/post/488752/


## Что используется
 * Golang для backend
 * Websocket для подключения игроков
 * Клиент написан на JS, используется фреймворк Phaser.js 
 * Обмен данными через json
 * Logrus для логов

## Что реализовано на данный момент
 * Чанковая система с бесконечной генерацией мира 
 * Тестовая система авторизации, пока просто набросок, в дальнейшем будет использована база данных
 * Сервер инициализации, где игрок получает данные для начала игры
 * Перемещение, простой алгоритм поиска пути
 * Генерация игровых объектов, пока что только деревьев
 * Сортировка в глубину для корректного отображения игровых объектов друг за другом
 * Управление мышью 



## Запуск
Для запуска сервера достаточно скопировать репозиторий на свою машину и вызвать go run:

```$ git clone https://github.com/Lolodin/MMORPG```       
```$ go run main.go```

Переходим в браузере по адресу localhost:8080

Первый прототип на чистом жс: https://lolodin.github.io/My2DGame/
