package main

import (
	"database/sql"
	"log"
	"os"

	configs "github.com/MatthewAraujo/airCast/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Importa suporte para leitura de arquivos de migração
	_ "github.com/jackc/pgx/v5/stdlib"                   // Driver PostgreSQL para database/sql
)

func main() {
	url := configs.Envs.Postgres.URL
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	// Configurando o driver de migração para PostgreSQL
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL driver: %v", err)
	}

	// Caminho para as migrações
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"pgx",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Verificando versão atual e estado do banco de dados
	v, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get version info: %v", err)
	}
	log.Printf("Version: %d, dirty: %v", v, dirty)

	// Forçando o estado limpo se o banco de dados estiver "sujo"
	if dirty {
		log.Println("Database is dirty. Forcing clean state...")
		if err := m.Force(int(v)); err != nil {
			log.Fatalf("Failed to force clean state: %v", err)
		}
	}

	// Executando comandos de migração
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
	}

	log.Println("Migration process completed successfully!")

}
