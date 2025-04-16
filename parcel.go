package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}

	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return p, err
	}
	defer db.Close()

	// заполните объект Parcel данными из таблицы
	rows, err := db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE number")
	if err != nil {
		return p, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return p, err
		}

	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	var res []Parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return res, err
	}
	defer db.Close()
	// заполните срез Parcel данными из таблицы
	rows, err := db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client ", sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt, &p.Client)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}
	statusCheck := Parcel{}
	err = rows.Scan(&statusCheck.Status)
	if err != nil {
		return err
	}

	if statusCheck.Status == "registered" {

		_, err = db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))

		if err != nil {
			return err
		}

		return nil
	} else {
		fmt.Println("Неверный статус")
		return nil
	}
}
func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}
	statusCheck := Parcel{}
	for rows.Next() {
		err = rows.Scan(&statusCheck.Status)
		if err != nil {
			return err
		}
	}

	if statusCheck.Status == "registered" {
		_, err := db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Неверный статус")
		return nil
	}

	return nil
}
