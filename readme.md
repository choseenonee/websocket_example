# Simple chat room realisation using websockets

История чата хранится в Postgres, если соединение с postgres пропало, отправляемые сообщения будут записываться в 
info.log файл.
<br>
Если нет соединения с postgres, не будет возможности получить сообщения, а также присоединиться к чату, даже имея его ID.