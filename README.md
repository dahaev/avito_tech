
Создать пользователя: POST: 127.0.0.1:8080/clients/

{
    "clientID":"CL902",
    "balanceMain":890.23
}

Возврат будет: 

{
    "clientID":"CL902",
    "balanceMain":890.23
}

Баланс в резерве:0(Значение дефолтное)

Создать заказ: POST 127.0.0.1:8080/orders/

{
    "OrderId":"ORD0001",
    "clientID":"CL902",
    "serviceID":"PR8003",
    "cost":350.68
}
ответ будет:

{
    "OrderId": "ORD0001",
    "clientID": "CL902",
    "serviceID": "PR8003",
    "cost": 350.68,
    "completed": false
}

В тоже самое время если запросить баланс пользователя CL902, по методу GET 127.0.0.1:8080/clients/CL902, ответ будет следующим: 

{
    "clientID": "CL902",
    "balanceMain": 539.55,
    "balanceReserved": 350.68
}

Прошу заметить, что типы услуг содержатся в отдельной таблице. Нельзя передать любой ID, сделал перечень в отдельной таблице.

Списание услуги. Направляем PATCH запрос, для перевода completed в true.

PATCH запрос выглядит следующим образом: 

127.0.0.1:8080/orders/ORD0001

Ответ:
{
    "OrderId": "ORD0001",
    "clientID": "CL902",
    "serviceID": "PR8003",
    "cost": 350.68,
    "completed": true
}

Проверяем баланс пользователя по методу GET 127.0.0.1:8080/clients/CL902:

{
    "clientID": "CL902",
    "balanceMain": 539.55,
    "balanceReserved": 0
}

Данные о выполненном заказе попадают в таблицу для бухгалетрии: 
[
    {
        "id": 12,
        "client_id": "CL902",
        "reportServiceID": "PR8003",
        "reportCost": 350.68
    }
]



