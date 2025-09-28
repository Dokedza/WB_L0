package migrator

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pressly/goose"
)

// для управления миграциями
type Migrator struct {
	db      *sql.DB
	timeout time.Duration
}

// новый мигратор
func New(db *sql.DB) *Migrator {
	return &Migrator{
		db:      db,
		timeout: 40 * time.Second,
	}
}

func (m *Migrator) SetTimeout(t time.Duration) {
	m.timeout = t
}

func (m *Migrator) WaitingDB() error {
	deadpool := time.Now().Add(m.timeout)
	log.Printf("Ожидание подключения базы данных (таймаут: %v)..", m.timeout)
	for time.Now().Before(deadpool) {
		err := m.db.Ping()
		if err == nil {
			log.Print("Подключение успешно")
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("Не удалось подклучиться к базе данных за %v", m.timeout)
}
func (m *Migrator) RunMigrations() error {
	err := m.WaitingDB()
	if err != nil {
		return fmt.Errorf("Ошибка ожидания базы данных: %w", err)
	}
	vers, err := goose.GetDBVersion(m.db)
	if err != nil {
		log.Printf("Миграции не применились ошибка: %w", err)
	} else {
		log.Printf("Текущая версия миграции: %d", vers)
	}

	log.Println("Начало применения миграций")
	err = goose.Up(m.db, "../migrations")
	if err != nil {
		return fmt.Errorf("Ошибка применения миграций: %w", err)
	}
	finalVers, err := goose.GetDBVersion(m.db)
	if err != nil {
		return fmt.Errorf("Ошибка получения финальной версии: %w", err)
	}
	log.Printf("Миграции выполнены успешно. Текущая версия: %d", finalVers)
	return nil
}
