# daec-ylyceum
Финальная задача 2 спринта. Яндекс Лицей

## Запуск 
Для запуска приложения требуется запущенный оркестратор и агент

### BASH (Linux/Windows with bash/Other UNIX-like) (Приоритетно)
```shell
./run.sh
```
### Windows (Командная строка)
```batch
./run.bat
```
### Manual
При ручном запуске переменные окружения нужно устанавливать вручную. Если не устанавливать то будут использованы стандартные. Что-бы запустить нужно скомпилировать cmd/main.go в server и agent соответственно и запустить скомпилированные файлы


##  Переменные окружения(конфигурация)

Задаются в файле .config

```
TIME_DIVISION_MS = <время деления в миллисекундах>
TIME_MULTIPLICATION_MS = <время умножения в миллисекундах>
TIME_SUBTRACTION_MS = <время вычитания в миллисекундах>
TIME_ADDITION_MS = <время сложения в миллисекундах>
COMPUTING_POWER = <количество запускаемых вычислителей>
SERVER_PORT = <порт открытия сервера>
```
## Примеры использования

### Веб-интерфейс
Открыт на пути localhost:<порт>/

### Запросы в API
Создание выражения для вычисления:
```shell
curl -d '{"expression": "2+2*2"}' localhost:<порт>/api/v1/calculate
```

Получение выражения по идентификатору:
```shell
curl localhost:<порт>/api/v1/expressions/<id>
```

Получение списка выражений:
```shell
curl localhost:<порт>/api/v1/expressions
```
## Как работает

1. Дергается ручка expression
2. Из выражения в виде строки состовляется дерево выражаний при помощи библотеки math-engine
  ![Image alt](https://study-and-dev.com/blog/contents/data/mediawiki/images/28/sda_book_22.png)
3. Запускается рекурсивный проход по дереву

Тут остановимся.

При проходе по дереву есть 2 типа узлов - операция (+-*/) или число. Так вот все эти узлы мы представляем как отдельные задачи. Для тех узлов что являются числами мы создаем задачу не требующую решения (потому что выражение из 1 числа уже посчитано) и выдаем ей идентификатор.
В узлах что считаются операциями мы устанавливаем идентификатор левой задачи и правой задачи (будь то число или другая операция). Это сделано для комбинирования просто чисел и результатов предыдуыщих выражений требуемых для вычисления.

Пример: 2+2\*2

Главной (последней) задачей будет умножение 2 на результат 2+2

Будут созданы задачи:

1. ИД: 1, Левая задача: 2, Правая задача: 5, Статус: Не посчитано, Операция: * 
2. ИД: 2, Левая задача: 3, Правая задача: 4, Статус: Не посчитано, Операция: +
3. ИД: 3, Левая задача: -1, Правая задача: -1, Статус: Посчитано, Результат=2
4. ИД: 4, Левая задача: -1, Правая задача: -1, Статус: Посчитано, Результат=2
5. ИД: 5, Левая задача: -1, Правая задача: -1, Статус: Посчитано, Результат=2

У нас получилось 5 задач (3 из них просто числа не требующие решения, 2 других операции требующие подсчета)

4. Агент идет в ручку /internal/task

Оркестратор проходит по всем задачам и ищет первую готовую для вычисления (та у которой левая и правая задача уже решены)

5. Выражение считается решенным когда задача являющаяся головой будет решена

### Было использовано
Бэкенд:
Logrus
math-engine
chi

Фронтенд:
Alpine.js
HTMX
Bootstrap





