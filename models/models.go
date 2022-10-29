package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	Client_id        string  `json:"clientID"`
	Balance_main     float64 `json:"balanceMain"`
	Balance_reserved float64 `json:"balanceReserved"`
}

type Service struct {
	Service_id string  `json:"Service_id"`
	cost       float64 `json:"cost"`
}

type Order struct {
	Order_id         string  `json:"OrderId"`
	Order_Client     string  `json:"clientID"`
	Order_service_id string  `json:"serviceID"`
	Cost             float64 `json:"cost"`
	Completed        bool    `json:"completed"`
}

type Report struct {
	ReportId        int     `json:"id"`
	ReportClientID  string  `json:"client_id"`
	ReportServiceID string  `json:"reportServiceID"`
	ReportCost      float64 `json:"reportCost"`
}

func GetClients() []Client {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println("Connect is OK")

	defer db.Close()
	result, err := db.Query("SELECT * FROM clients")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	clients := []Client{}
	for result.Next() {
		var cli Client
		err := result.Scan(&cli.Client_id, &cli.Balance_main, &cli.Balance_reserved)
		if err != nil {
			panic(err)
		}
		clients = append(clients, cli)
	}
	fmt.Println(clients)
	return clients
}

func AddClient(client Client) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println("Connect is OK")
	defer db.Close()
	if client.Balance_main < 0 {
		return
	}
	insert, err := db.Query("INSERT INTO clients (client_id, balance_main) VALUES (?, ?)", client.Client_id, client.Balance_main)
	if err != nil {
		panic(err)
	}
	defer insert.Close()
}

func AddOrder(order Order) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println("Connect is ok")
	defer db.Close()
	db_client, err := db.Query("SELECT * FROM clients")
	if err != nil {
		panic(err)
	}
	for db_client.Next() {
		var client Client
		err := db_client.Scan(&client.Client_id, &client.Balance_main, &client.Balance_reserved)
		if err != nil {
			panic(err)
		}
		fmt.Println(client)
		fmt.Println("Зашел в db_client")
		if order.Order_Client == client.Client_id {
			fmt.Println("Зашел в IF")
			new_balance := client.Balance_main - order.Cost
			new_resrv := client.Balance_reserved + order.Cost
			db_update, err := db.Query("UPDATE clients SET balance_main=?, balance_reserved=? WHERE client_id=?", new_balance, new_resrv, order.Order_Client)
			defer db_update.Close()
			if err != nil {
				panic(err)
			}
		}
	}
	insert, err := db.Query("INSERT INTO orders (order_id, client, service_id, cost) VALUES (?, ?, ?, ?)", order.Order_id, order.Order_Client, order.Order_service_id, order.Cost)
	if err != nil {
		panic(err)
	}
	defer insert.Close()
}

func FinanceReport(order Order) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	db_order, err := db.Query("SELECT * FROM clients")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	for db_order.Next() {
		var cli Client
		err := db_order.Scan(&cli.Client_id, &cli.Balance_main, &cli.Balance_reserved)
		if err != nil {
			fmt.Println(err)
		}
		if cli.Client_id == order.Order_id {
			insert, err := db.Query("UPDATE clients SET balance_reserved=? WHERE client_id=?", cli.Balance_reserved-order.Cost)
			if err != nil {
				fmt.Println(err)
			}
			insert.Close()
		}
	}

}

func GetOrderById(code string) *Order {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	ord := &Order{}
	fmt.Println("Начальный вовод ord", ord)
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println("Подключение к DB успешно")
	defer db.Close()
	result, err := db.Query("SELECT * FROM orders WHERE order_id=?", code)
	if err != nil {
		fmt.Println("Query не сработал")
		fmt.Println("Err", err.Error())
		return nil
	}
	fmt.Println("code: ", code)
	fmt.Println("result: ", result)
	if result.Next() {
		fmt.Println(" цикл зашли")
		err = result.Scan(&ord.Order_id, &ord.Order_Client, &ord.Order_service_id, &ord.Cost, &ord.Completed)
		if err != nil {
			fmt.Println("Сработал if")
			return nil
		}
	}
	insert_rep, err := db.Query("INSERT INTO report(client_id, service_id, cost) VALUES (?,?,?)", ord.Order_Client, ord.Order_service_id, ord.Cost)
	fmt.Println("INSERT в таблицу отчетов")
	defer insert_rep.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("under insert")
	insert_client, err := db.Query("SELECT * FROM clients")
	fmt.Println("after insert")
	if err != nil {
		panic(err)
	}
	defer insert_client.Close()
	for insert_client.Next() {
		fmt.Println("Мы в цикле обновления таблицы clients")
		var cli Client
		err := insert_client.Scan(&cli.Client_id, &cli.Balance_main, &cli.Balance_reserved)
		if err != nil {
			panic(err)
		}
		fmt.Println(cli.Client_id, code)
		if cli.Client_id == ord.Order_Client {
			fmt.Println("Зашел в обновление таблицы clients")
			cli.Balance_reserved = cli.Balance_reserved - ord.Cost
			db_update, err := db.Query("UPDATE clients SET balance_reserved=? WHERE client_id=?", cli.Balance_reserved, ord.Order_Client)
			if err != nil {
				panic(err)
			}
			defer db_update.Close()

		}
	}
	db_order_del, err := db.Query("DELETE FROM orders WHERE order_id=?", code)
	if err != nil {
		fmt.Println(err)
	}
	db_order_del.Close()
	fmt.Println("Вывод ORD", ord)

	return ord
}

func GetClientByID(client_id string) *Client {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	cli := &Client{}
	if err != nil {
		fmt.Println("err", err.Error())
		return nil
	}
	defer db.Close()
	results, err := db.Query("SELECT * FROM clients WHERE client_id=?", client_id)
	if err != nil {
		fmt.Println("err", err.Error())
	}
	if results.Next() {
		err = results.Scan(&cli.Client_id, &cli.Balance_main, &cli.Balance_reserved)
		if err != nil {
			return nil
		}
	}
	return cli
}

func GetReports() []Report {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/avito_tech")
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println("Connect is OK")

	defer db.Close()
	result, err := db.Query("SELECT * FROM report")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	report := []Report{}
	for result.Next() {
		var rep Report
		err := result.Scan(&rep.ReportId, &rep.ReportClientID, &rep.ReportServiceID, &rep.ReportCost)
		if err != nil {
			panic(err)
		}
		report = append(report, rep)
	}
	fmt.Println(report)
	return report
}
