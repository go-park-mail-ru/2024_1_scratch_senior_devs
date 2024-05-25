# Нагрузочное тестирование

## Создание заметок
wrk -t12 -c400 -d30s --script=script.lua --latency http://127.0.0.1:8080
![img.png](images/img.png)

## Чтение заметок
wrk -t12 -c400 -d30s --script=script2.lua --latency http://127.0.0.1:8080
![img_1.png](images/img_1.png)
