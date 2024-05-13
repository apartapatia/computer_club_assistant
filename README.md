# Ассистент компьютерного клуба

Прототип системы, которая следит за работой компьютерного клуба, обрабатывает события и подсчитывает выручку за день и время занятости каждого стола.


## Формат входных данных
    <количество столов в компьютерном клубе>
    <время начала работы> <время окончания работы>
    <стоимость часа в компьютерном клубе>
    <время события 1> <идентификатор события 1> <тело события 1>
    <время события 2> <идентификатор события 2> <тело события 2>
                            ...
    <время события N> <идентификатор события N> <тело события N>

Входные данные задаются файлом в формате `.txt`, они должны находиться в директории `/configs`.

## Запуск приложения

Склонируйте репозиторий и перейдите в корневую папку проекта.

### Запуск контейнера 🐋 docker
Для запуска контейнера, укажите имя файла в переменной окружения `FILE_NAME`:

```
docker build -t computer_club_assistant:v1 .

docker run -e FILE_NAME=test_main.txt computer_club_assistant:v1
```
Где `test_main.txt` - это имя вашего файла с данными. Убедитесь, что файл находится в директории `/configs`.

### Запуск на 🪟 Windows

```
go build -o computer_club_assistant.exe cmd/computer_club_assistant/main.go

./computer_club_assistant.exe <file_name>
```

### Запуск на 🐧 Linux

```
go build -o computer_club_assistant cmd/computer_club_assistant/main.go

./computer_club_assistant <file_name>
```
