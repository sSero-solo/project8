package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	result, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add parcel: %w", err)
	}

	// верните идентификатор последней добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := s.db.QueryRow(
		"SELECT number, client, status, address, created_at FROM parcel WHERE number = ?",
		number,
	).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, fmt.Errorf("parcel with number %d not found", number)
		}
		return Parcel{}, fmt.Errorf("failed to get parcel: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = ?",
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get parcels by client: %w", err)
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parcel: %w", err)
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	result, err := s.db.Exec(
		"UPDATE parcel SET status = ? WHERE number = ?",
		status, number,
	)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	// Проверяем текущий статус
	var currentStatus string
	err := s.db.QueryRow(
		"SELECT status FROM parcel WHERE number = ?",
		number,
	).Scan(&currentStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel with number %d not found", number)
		}
		return fmt.Errorf("failed to check parcel status: %w", err)
	}

	if currentStatus != ParcelStatusRegistered {
		return errors.New("can only change address for registered parcels")
	}

	// Обновляем адрес
	result, err := s.db.Exec(
		"UPDATE parcel SET address = ? WHERE number = ?",
		address, number,
	)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	// Проверяем текущий статус
	var currentStatus string
	err := s.db.QueryRow(
		"SELECT status FROM parcel WHERE number = ?",
		number,
	).Scan(&currentStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel with number %d not found", number)
		}
		return fmt.Errorf("failed to check parcel status: %w", err)
	}

	if currentStatus != ParcelStatusRegistered {
		return errors.New("can only delete registered parcels")
	}

	// Удаляем посылку
	result, err := s.db.Exec(
		"DELETE FROM parcel WHERE number = ?",
		number,
	)
	if err != nil {
		return fmt.Errorf("failed to delete parcel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}
