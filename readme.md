# Simple chat room realisation using websockets

История чата хранится в Postgres, если соединение с postgres пропало, отправляемые сообщения будут записываться в 
info.log файл.
<br>
Если нет соединения с postgres, не будет возможности получить сообщения, а также присоединиться к чату, даже имея его ID.
## Логическая схема
![schema.jpg](schema.jpg)
## TODOs
- [ ] Benchmark 
    - [ ] Нагрузочное тестирование
        - [x] Создание N чатов
        - [x] Создание N горутин-слушателей для каждого чата
        - [x] Отправка N сообщений подряд из 1 горутины
        - [ ] Отправка N сообщений в K времени (возможно, из разных горутин)
- [ ] Пул воркеров, изменяющийся под нагрузкой
    - [ ] Структуру-мастер, управляющую воркерами
    - [ ] Функция воркер
      - [ ] Воркер обслуживает сообщения чата?
      - [ ] Воркер обслуживает 1 сообщение для 1 получателя?
