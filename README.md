1 - Enpdoint /sendMoney
    1.1 - Transaction: id / from / to
    1.2 - 
2 - Транзакция совершена => записали в базу + посылаем событие в Kafka. Из Kafka читает консьюмер, который и отправляет данные (берет из базы) в налоговую.
3 - Handler
    3.1 - usecase