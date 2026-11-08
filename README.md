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
```json
{
  "started_at": "2026-11-08T23:50:06.893416125+03:00",
  "finished_at": "2026-11-08T23:50:06.90309248+03:00",
  "targets_count": 1,
  "tasks_count": 10,
  "findings": [
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%D0%9F%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B0%D0%B2%D0%BB%D1%8F%D1%82%D1%8C+%D0%BD%D0%B0+%D1%80%D0%B0%D1%81%D1%81%D0%BC%D0%BE%D1%82%D1%80%D0%B5%D0%BD%D0%B8%D0%B5\u0026id=1%22",
      "parameter": "id",
      "payload": "1\"",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893598596+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=1%22\u0026id=1",
      "parameter": "Submit",
      "payload": "1\"",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893679373+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%D0%9F%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B0%D0%B2%D0%BB%D1%8F%D1%82%D1%8C+%D0%BD%D0%B0+%D1%80%D0%B0%D1%81%D1%81%D0%BC%D0%BE%D1%82%D1%80%D0%B5%D0%BD%D0%B8%D0%B5\u0026id=1%27+UNION+SELECT+NULL+--",
      "parameter": "id",
      "payload": "1' UNION SELECT NULL --",
      "vulnerable": true,
      "evidence": "\u003cpre\u003eYou have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right",
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893593111+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=1%27\u0026id=1",
      "parameter": "Submit",
      "payload": "1'",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893723177+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%D0%9F%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B0%D0%B2%D0%BB%D1%8F%D1%82%D1%8C+%D0%BD%D0%B0+%D1%80%D0%B0%D1%81%D1%81%D0%BC%D0%BE%D1%82%D1%80%D0%B5%D0%BD%D0%B8%D0%B5\u0026id=1%27",
      "parameter": "id",
      "payload": "1'",
      "vulnerable": true,
      "evidence": "\u003cpre\u003eYou have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right",
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893639691+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%D0%9F%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B0%D0%B2%D0%BB%D1%8F%D1%82%D1%8C+%D0%BD%D0%B0+%D1%80%D0%B0%D1%81%D1%81%D0%BC%D0%BE%D1%82%D1%80%D0%B5%D0%BD%D0%B8%D0%B5\u0026id=1%27+OR+%271%27%3D%271%27+--",
      "parameter": "id",
      "payload": "1' OR '1'='1' --",
      "vulnerable": true,
      "evidence": "\u003cpre\u003eYou have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right",
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893661724+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=1%27+OR+%271%27%3D%271%27+--\u0026id=1",
      "parameter": "Submit",
      "payload": "1' OR '1'='1' --",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.89369151+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%27\u0026id=1",
      "parameter": "Submit",
      "payload": "'",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893701689+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=1%27+UNION+SELECT+NULL+--\u0026id=1",
      "parameter": "Submit",
      "payload": "1' UNION SELECT NULL --",
      "vulnerable": false,
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893558952+03:00"
    },
    {
      "url": "http://localhost:8080/vulnerabilities/sqli/?Submit=%D0%9F%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B0%D0%B2%D0%BB%D1%8F%D1%82%D1%8C+%D0%BD%D0%B0+%D1%80%D0%B0%D1%81%D1%81%D0%BC%D0%BE%D1%82%D1%80%D0%B5%D0%BD%D0%B8%D0%B5\u0026id=%27",
      "parameter": "id",
      "payload": "'",
      "vulnerable": true,
      "evidence": "\u003cpre\u003eYou have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right",
      "status_code": 200,
      "checked_at": "2026-11-08T23:50:06.893486258+03:00"
    }
  ]
}
```

## Known Limitations
- **DVWA High:** false negative. Уровень High читает payload из cookie `id`, а не из GET-параметра. Тул проверяет только GET-каналы.
- **Error-based only:** детекция работает только на SQL-ошибках в теле ответа. Blind SQLi (boolean/time-based) не поддерживается.
- **GET only:** POST-параметры, заголовки, cookie как каналы ввода — в бэклоге.

## Архитектура
cmd/sqli-detector/main.go    — точка входа, связка компонентов
internal/config/             — парсинг и валидация флагов
internal/scanner/            — воркер-пул, HTTP-запросы, детекция
internal/report/             — JSON-сериализация отчёта
