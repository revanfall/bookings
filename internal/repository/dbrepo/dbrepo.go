package dbrepo

import (
	"database/sql"
	"github.com/revanfall/bookings/internal/config"
	"github.com/revanfall/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{a, conn}
}
