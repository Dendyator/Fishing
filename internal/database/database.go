package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Подключение к базе данных успешно")
	return &DB{db}, nil
}

// CreateFishingData создает новую запись о рыбалке.
func (db *DB) CreateFishingData(data FishingData) error {
	query := `
        INSERT INTO fishing_data (date, location, coordinates, comment, fishing_methods, caught_fish, trophy_fish)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	caughtFishJSON, err := json.Marshal(data.CaughtFish)
	if err != nil {
		return fmt.Errorf("ошибка преобразования пойманной рыбы в JSON: %v", err)
	}

	trophyFishJSON, err := json.Marshal(data.TrophyFish)
	if err != nil {
		return fmt.Errorf("ошибка преобразования трофейной рыбы в JSON: %v", err)
	}

	_, err = db.Exec(query,
		data.Date,
		data.Location,
		data.Coordinates,
		data.Comment,
		pq.Array(data.FishingMethods),
		string(caughtFishJSON),
		string(trophyFishJSON),
	)
	return err
}

// GetFishingDataByID получает данные о рыбалке по ID.
func (db *DB) GetFishingDataByID(id int) (*FishingData, error) {
	query := `
        SELECT id, date, location, coordinates, comment, fishing_methods, caught_fish, trophy_fish
        FROM fishing_data
        WHERE id = $1
    `

	var data FishingData
	var caughtFishJSON, trophyFishJSON sql.NullString
	err := db.QueryRow(query, id).Scan(
		&data.ID,
		&data.Date,
		&data.Location,
		&data.Coordinates,
		&data.Comment,
		pq.Array(&data.FishingMethods),
		&caughtFishJSON,
		&trophyFishJSON,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("нет данных о рыбалке для указанного ID")
	} else if err != nil {
		return nil, fmt.Errorf("ошибка получения данных о рыбалке: %v", err)
	}

	if caughtFishJSON.Valid {
		err = json.Unmarshal([]byte(caughtFishJSON.String), &data.CaughtFish)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования пойманной рыбы из JSON: %v", err)
		}
	}

	if trophyFishJSON.Valid {
		var trophy FishingDataTrophy
		err = json.Unmarshal([]byte(trophyFishJSON.String), &trophy)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования трофейной рыбы из JSON: %v", err)
		}
		data.TrophyFish = &trophy
	}

	return &data, nil
}

// GetAllFishingData получает все записи о рыбалке.
func (db *DB) GetAllFishingData() ([]FishingData, error) {
	query := `
        SELECT id, date, location, coordinates, comment, fishing_methods, caught_fish, trophy_fish
        FROM fishing_data
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных о рыбалке: %v", err)
	}
	defer rows.Close()

	var data []FishingData
	for rows.Next() {
		var record FishingData
		var caughtFishJSON, trophyFishJSON sql.NullString
		err := rows.Scan(
			&record.ID,
			&record.Date,
			&record.Location,
			&record.Coordinates,
			&record.Comment,
			pq.Array(&record.FishingMethods),
			&caughtFishJSON,
			&trophyFishJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}

		if caughtFishJSON.Valid {
			err = json.Unmarshal([]byte(caughtFishJSON.String), &record.CaughtFish)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования пойманной рыбы из JSON: %v", err)
			}
		}

		if trophyFishJSON.Valid {
			var trophy FishingDataTrophy
			err = json.Unmarshal([]byte(trophyFishJSON.String), &trophy)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования трофейной рыбы из JSON: %v", err)
			}
			record.TrophyFish = &trophy
		}

		data = append(data, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %v", err)
	}

	return data, nil
}

// UpdateFishingData обновляет существующую запись о рыбалке.
func (db *DB) UpdateFishingData(id int, data FishingData) error {
	query := `
        UPDATE fishing_data
        SET date = $1, location = $2, coordinates = $3, comment = $4, fishing_methods = $5, caught_fish = $6, trophy_fish = $7
        WHERE id = $8
    `

	caughtFishJSON, err := json.Marshal(data.CaughtFish)
	if err != nil {
		return fmt.Errorf("ошибка преобразования пойманной рыбы в JSON: %v", err)
	}

	trophyFishJSON, err := json.Marshal(data.TrophyFish)
	if err != nil {
		return fmt.Errorf("ошибка преобразования трофейной рыбы в JSON: %v", err)
	}

	result, err := db.Exec(query,
		data.Date,
		data.Location,
		data.Coordinates,
		data.Comment,
		pq.Array(data.FishingMethods),
		string(caughtFishJSON),
		string(trophyFishJSON),
		id,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления данных о рыбалке: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества затронутых строк: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("не найдена запись о рыбалке с указанным ID")
	}

	return nil
}

// GetAllFishingMethods получает список всех способов ловли.
func (db *DB) GetAllFishingMethods() ([]string, error) {
	query := "SELECT name FROM fishing_methods"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения способов ловли: %v", err)
	}
	defer rows.Close()

	var methods []string
	for rows.Next() {
		var method string
		if err := rows.Scan(&method); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		methods = append(methods, method)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %v", err)
	}

	return methods, nil
}

// InitializeData инициализирует базовые данные в базе данных.
func (db *DB) InitializeData(fishingMethods, fishSpecies, baits []string) error {
	// Добавление способов ловли
	for _, method := range fishingMethods {
		if err := db.CreateFishingMethod(method); err != nil && !isDuplicateError(err) {
			return fmt.Errorf("не удалось создать способ ловли %s: %v", method, err)
		}
	}

	// Добавление видов рыб
	for _, specie := range fishSpecies {
		if err := db.CreateFishSpecies(specie); err != nil && !isDuplicateError(err) {
			return fmt.Errorf("не удалось создать вид рыбы %s: %v", specie, err)
		}
	}

	// Добавление приманок
	for _, bait := range baits {
		if err := db.CreateBait(bait); err != nil && !isDuplicateError(err) {
			return fmt.Errorf("не удалось создать приманку %s: %v", bait, err)
		}
	}

	return nil
}

// Вспомогательная функция для проверки дубликатов
func isDuplicateError(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		return err.Code.Name() == "unique_violation"
	}
	return false
}

// CreateFishingMethod создает новый способ ловли.
func (db *DB) CreateFishingMethod(name string) error {
	query := "INSERT INTO fishing_methods (name) VALUES ($1)"
	_, err := db.Exec(query, name)
	return err
}

// CreateFishSpecies создает новый вид рыбы.
func (db *DB) CreateFishSpecies(name string) error {
	query := "INSERT INTO fish_species (name) VALUES ($1)"
	_, err := db.Exec(query, name)
	return err
}

// CreateBait создает новую приманку.
func (db *DB) CreateBait(name string) error {
	query := "INSERT INTO baits (name) VALUES ($1)"
	_, err := db.Exec(query, name)
	return err
}

func (db *DB) GetAllFishSpecies() ([]string, error) {
	query := "SELECT name FROM fish_species"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения видов рыб: %v", err)
	}
	defer rows.Close()

	var species []string
	for rows.Next() {
		var specie string
		if err := rows.Scan(&specie); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		species = append(species, specie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %v", err)
	}

	return species, nil
}

func (db *DB) GetAllBaits() ([]string, error) {
	query := "SELECT name FROM baits"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приманок: %v", err)
	}
	defer rows.Close()

	var baits []string
	for rows.Next() {
		var bait string
		if err := rows.Scan(&bait); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		baits = append(baits, bait)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %v", err)
	}

	return baits, nil
}
