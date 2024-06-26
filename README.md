# Arithmetic Operations


<details><summary><b>Требования к заданию (часть первая)</b></summary>

Пользователь хочет считать арифметические выражения. 
Он вводит строку `2 + 2 * 2` и хочет получить в ответ `6`. 
Но наши операции сложения и умножения (также деления и вычитания) выполняются **"очень-очень" долго**. 
Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, **невозможна**. 
Более того: вычисление каждой такой операции в нашей **"альтернативной реальности"** занимает **"гигантские"** вычислительные мощности. 
Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "**машин**". 
Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? 
Если выражение наконец будет вычислено - то он получит результат. 
Помните, что некоторые части арфиметического выражения можно вычислять **параллельно**.

## Front-end часть

### GUI, который можно представить как 4 страницы

1) Форма ввода арифметического выражения. Пользователь вводит арифметическое выражение и отправляет **POST** http-запрос с этим выражением на back-end. Примечание: Запросы должны быть **идемпотентными**. К запросам добавляется **уникальный идентификатор**. Если пользователь отправляет запрос с идентификатором, который уже отправлялся и был принят к обработке - ответ 200. Возможные варианты ответа:
    - _200_ - Выражение успешно принято, распаршено и принято к обработке
    - _400_ - Выражение невалидно
    - _500_ - Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
2) Страница со списком выражений в виде списка с выражениями. Каждая запись на странице содержит статус, выражение, дату его создания и дату заверщения вычисления. Страница получает данные GET http-запрсом с back-end-а
3) Страница со списком операций в виде пар: имя операции + время его выполнения (доступное для редактирования поле). Как уже оговаривалось в условии задачи, наши операции выполняются "как будто бы очень долго". Страница получает данные GET http-запрсом с back-end-а. Пользователь может настроить время выполения операции и сохранить изменения.
4) Страница со списком вычислительных можностей. Страница получает данные GET http-запросом с сервера в виде пар: имя вычислительного ресурса + выполняемая на нём операция.

### Требования:

1) Оркестратор может перезапускаться без потери состояния. Все выражения храним в СУБД.
2) Оркестратор должен отслеживать задачи, которые выполняются слишком долго (вычислитель тоже может уйти со связи) и делать их повторно доступными для вычислений.


## Back-end часть

### Состоит из 2 элементов:

- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
- Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.

### Оркестратор
Сервер, который имеет следующие endpoint-ы:

- Добавление вычисления арифметического выражения.
- Получение списка выражений со статусами.
- Получение значения выражения по его идентификатору.
- Получение списка доступных операций со временем их выполения.
- Получение задачи для выполения.
- Приём результата обработки данных.


### Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения. При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды.


</details>

<details><summary><b>Требования к заданию (часть вторая)</b></summary>
Продолжаем работу над проектом `Распределенный калькулятор`.

В этой части работы над проектом реализуем **персистентность** и **многопользовательский режим**.

### Функционал:
1. Добавляем регистрацию пользователя. В ответ получает 200 в случае успеха. В противном случае - 401.
    ```http request
    POST /signup
   
    {
        "login": "login_value",
        "password": "password_value"
    }
    ```

2. Добавляем вход. В ответ получает 200 и JWT токен для последующей авторизации. В противном случае - 403.
    ```http request
    POST /login
    
    {
        "login": "login_value",
        "password": "password_value"
    }
    ```

### Баллы:
1. Весь реализованный ранее функционал работает как раньше, только в контексте конкретного пользователя.
   - 20 баллов.
2. У кого выражения хранились в памяти - переводим хранение в СУБД.
   - 20 баллов
3. У кого общение вычислителя и сервера вычислений было реализовано с помощью HTTP - переводим взаимодействие на GRPC.
   - 10 баллов
4. Покрытие проекта модульными тестами
   - 10 баллов
5. Покрытие проекта интеграционными тестами
   - 10 баллов

### Правила оформления:
- проект находится на GitHub
- к проекту прилагается файл с подробным описанием (как запустить, проверить функционал и протестировать)
- отдельным блоком идут подробно описанные тестовые сценарии
- автоматизируйте поднятие окружения для запуска вашей программы

</details>

http://localhost:8080/

## По поводу того что есть по критериям:
- Весь реализованный ранее функционал работает как раньше, только в контексте конкретного пользователя   &#10004;
- У кого выражения хранились в памяти - переводим хранение в СУБД   &#10004;
- У кого общение вычислителя и сервера вычислений было реализовано с помощью HTTP - переводим взаимодействие на GRPC &#x2717;
- Покрытие проекта модульными тестами &#10004;
- Покрытие проекта интеграционными тестами &#10004;
## Инструкция к запуску:
### Способ с Docker:
1. Скачать [Docker](https://www.docker.com/products/docker-desktop/)
2. Склонировать проект или скачать `git clone github.com/byoverr/arithmetic_operations`
3. Перейти в папку проекта с помощью `cd <путь к файлу>`
4. Ввести команду `docker compose -f docker-compose.yml -p arithmetic_operations up`
5. Готово, все запустилось, по тестам читать далее

P.S Чтобы удалить контейнеры и образы используйте такие команды как:

`docker ps ; docker stop <ID контейнера> ; docker rm <ID контейнера>`

`docker images ; docker rmi <ID образа>`
### Способ дедовский:
1. Склонировать проект или скачать `git clone github.com/byoverr/arithmetic_operations`
2. Установить PostgreSQL
3. Поменять в config.json данные о PostgreSQL(и другие если хотите)
4. Установить все зависимости из go.mod `go mod download` `go mod tidy`
5. Включить сервер PostgreSQL
6. Запустить файл main.go

## Какие ошибки обрабатываются в выражении:

Сначала RemoveAllSpaces убирает пробелы в выражении
- HasDoubleSymbol проверяет на двойной символ
- IsValidParentheses проверяет скобочную последовательность
- HasDivizionByZero проверяет есть ли деление на ноль
- HasValidCharacters проверяет на допустимые символы
- ContainsCorrectFloatPoint проверяет на точку в правильном месте(должно быть в числе float)
- HasAtLeastOneExpression проверяет на хотя бы одно выражение число оператор число
- ExpressionStartsWithNumber проверяет - первым ли идёт число или скобка в выражении
## Как работает подсчёт выражения:

1. Делаем постфиксную запись выражения
2. Ищем выражения вида `24+` (число число оператор) или ищем выражения вида `2 + 2 +` (число оператор число оператор(операторы одинаковые))
3. Отправляем субвыражения на горутины(агенты)
4. Подставляем обратно в постфиксную запись наши субвыражения
5. Повторяем пока не останется одно число

## Какие запросы можем сделать и как:

**Есть файл [Postman](https://github.com/byoverr/arithmetic_operations/blob/main/docs/arithmetic_operations.postman_collection.json)**

Запросы делал через Postman
- `http://localhost:8080/register` - метод POST и Body такое(регистрация):

```json
{
   "username": "someuser@gmail.com",
   "password": "Qwerty_123"
}
```
- `http://localhost:8080/login` - метод POST и Body такое же как в register(вход)

Дает JWT токен, который впоследствии надо вписывать во все запросы(засовывать в Header так: "Authorization": "Bearer <JWT токен>"), написанные ниже
- `http://localhost:8080/expression` - метод POST и Body например такое(отправляем выражение на обработку):
```json
{
    "expression": "(2+2)*2"
}
```
![alt tag](https://github.com/byoverr/arithmetic_operations/blob/main/docs/img/example.png "Пример")
- `http://localhost:8080/expression` метод GET - получаем все выражения
- `http://localhost:8080/expression/1` метод GET - получаем выражение с id = 1
- `http://localhost:8080/operation` метод GET - получаем все операции(длительность операции)
- `http://localhost:8080/operation` метод PUT - обновляем операцию тут Body будет например такое:
```json
{
   "OperationKind": "addition",
   "durationInMilliSecond": 20
}
```
- `http://localhost:8080/agent` метод PUT - добавляем одного агента
- - `http://localhost:8080/agent` метод DELETE - удаляем одного агента
- `http://localhost:8080/agents` метод GET - получаем всех агентов
## Как это все работает:

![alt tag](https://github.com/byoverr/arithmetic_operations/blob/main/docs/img/scheme.png "Схемка")


У нас читается конфиг с файла config.json, инициализируется БД, chi роутер, выполняются непосчитанные выражения и привязываются паттерны. По вызванному паттерну есть хендлер в котором обрабатывается запрос и идет работа с БД, создаются таски с выражениями и длительностью операций. Всё обрабатывается и возвращается ответ.

Используемые технологии:

*chi, viper, pgx, slog, postgresql, docker(я пытался), JWT*

По всем вопросам и не вопросам пишите в [Telegram](https://t.me/super_serejka)

Special thanks to this [video](https://youtu.be/rCJvW2xgnk0?si=0bLCG5tMzKORbMxo)
