package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Sigue siendo el driver de MySQL compatible con MariaDB
	"time"
)

// Configuración de la base de datos
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Inicializar conexión a MySQL (compatible con MariaDB)
func NewMySQLConnection(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	// Intentar abrir la conexión
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar que la conexión funciona
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Inicializar las tablas de la base de datos si es necesario
	if err := initTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// Crear tablas si no existen
func initTables(db *sql.DB) error {
	// Tabla de usuarios
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			INDEX (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}

	// Tabla de datos de sensores
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_data (
			id INT AUTO_INCREMENT PRIMARY KEY,
			temperatura_dht FLOAT NOT NULL,
			luz FLOAT NOT NULL,
			humedad FLOAT NOT NULL,
			humo FLOAT NOT NULL,
			created_at DATETIME NOT NULL,
			INDEX (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}

	// Tabla de alertas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS alerts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			sensor_id INT NOT NULL,
			sensor_type VARCHAR(20) NOT NULL,
			value FLOAT NOT NULL,
			message TEXT NOT NULL,
			is_read BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL,
			INDEX (sensor_id),
			INDEX (sensor_type),
			INDEX (is_read),
			INDEX (created_at),
			FOREIGN KEY (sensor_id) REFERENCES sensor_data(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}

	return nil
}
