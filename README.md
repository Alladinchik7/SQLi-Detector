# sqli-detector

Error-based SQL-инъектор на Go с конкурентной обработкой URL.

## Возможности
- Многопоточное сканирование (-t)
- Автоматическая генерация задач по всем GET-параметрам
- Error-based детекция (MySQL, PostgreSQL, SQLite, Oracle)
- Экспорт в JSON с evidence-фрагментами

## Установка
go install github.com/Alladinchik/slqi-detector/cmd/sqli-detector@latest

## Использование
# Один URL
sqli-detector -u "http://target/page?id=1" -t 10

# Список URL из файла
sqli-detector -f urls.txt -o report.json

# С cookie (для аутентифицированных целей)
sqli-detector -u "http://target/page?id=1" -cookie "PHPSESSID=xxx; security=low"

## Флаги
-u string       один URL для проверки
-f string       файл со списком URL (по одному на строку)
-t int          число воркеров (default 10)
-o string       путь к JSON-отчёту (default stdout)
-cookie string  сырое значение заголовка Cookie
-timeout duration таймаут запроса (default 10s)

## Пример отчёта
<вставь сюда красивый JSON-фрагмент из low.json>

## Known Limitations
- **DVWA High:** false negative. Уровень High читает payload из cookie `id`, а не из GET-параметра. Тул проверяет только GET-каналы.
- **Error-based only:** детекция работает только на SQL-ошибках в теле ответа. Blind SQLi (boolean/time-based) не поддерживается.
- **GET only:** POST-параметры, заголовки, cookie как каналы ввода — в бэклоге.

## Архитектура
cmd/sqli-detector/main.go    — точка входа, связка компонентов
internal/config/             — парсинг и валидация флагов
internal/scanner/            — воркер-пул, HTTP-запросы, детекция
internal/report/             — JSON-сериализация отчёта