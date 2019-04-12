# go-sendMail
Формирование и отправка письма с содержимым указанной папки

Параметры хранятся в ini-файлах рядом с приложением. Количество ini-файлов не ограничено

**Параметры ini-файла**

[smtp]

server   = # Адрес SMTP сервера

port     = # Порт SMTP сервера

login    = # Логин

password = # Пароль

[mail]

from_addr = # Адрес отправителя (должен совпадать с логином и сервером отправки)

from_name = # Текстовая метка отправителя

to_addr   = # Адрес получателя

to_name   = # Текстовая метка получателя

subject   = # Заголовок письма, поддерживаются шаблоны "<date> <time>"

body      = # Тело письма (поддерживаются <date> <time>)

[paths]

7z           = # C:\Program Files\7-Zip\7z.exe

directory    = # Путь до папки, содержимое которой необходимо прислать

archive_name = # имя архива (поддерживаются "<date> <time>")

archive_type = # тип архива (7z или zip)
