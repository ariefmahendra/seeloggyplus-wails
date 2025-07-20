package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"seeloggyplus/backend/config"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
)

type ServerRepository interface {
	Create(ctx context.Context, server *entity.Server) (*entity.Server, error)
	Update(ctx context.Context, server *entity.Server) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Server, error)
	FindByID(ctx context.Context, id string) (*entity.Server, error)
	MigrateTable(ctx context.Context) error
}

type serverRepositoryImpl struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) ServerRepository {
	return &serverRepositoryImpl{db: db}
}

func (s *serverRepositoryImpl) MigrateTable(ctx context.Context) error {
	log := logger.Get()
	log.Debug().Msg("Executing migrate table query")

	_, err := s.db.ExecContext(ctx, config.ServerMigrateTable)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute migrate table query")
		return err
	}
	return nil
}

func (s *serverRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Server, error) {
	log := logger.Get()
	log.Debug().Str("id", id).Msg("Executing FindByID query")

	var serverById entity.Server
	err := s.db.QueryRowContext(ctx, config.GetServerById, id).
		Scan(&serverById.ID, &serverById.Name, &serverById.Address, &serverById.Port, &serverById.User, &serverById.Password, &serverById.CreatedAt, &serverById.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Str("id", id).Msg("Server not found in DB")
			return nil, nil
		}
		log.Error().Err(err).Str("id", id).Msg("Failed to scan server row")
		return nil, err
	}

	return &serverById, nil
}

func (s *serverRepositoryImpl) Create(ctx context.Context, server *entity.Server) (*entity.Server, error) {
	log := logger.Get()
	log.Debug().Str("name", server.Name).Msg("Executing Create query")

	var serverCreated entity.Server
	err := s.db.QueryRowContext(ctx, config.CreateServer, server.ID, server.Name, server.Address, server.Port, server.User, server.Password).
		Scan(&serverCreated.ID, &serverCreated.Name, &serverCreated.Address, &serverCreated.Port, &serverCreated.User, &serverCreated.Password, &serverCreated.CreatedAt, &serverCreated.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Str("name", server.Name).Msg("Failed to execute create query")
		return nil, err
	}

	return &serverCreated, nil
}

func (s *serverRepositoryImpl) Update(ctx context.Context, server *entity.Server) error {
	log := logger.Get()
	log.Debug().Str("id", server.ID).Msg("Executing Update query")

	res, err := s.db.ExecContext(ctx, config.UpdateServer, server.Name, server.Address, server.Port, server.User, server.Password, server.ID)
	if err != nil {
		log.Error().Err(err).Str("id", server.ID).Msg("Failed to execute update query")
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Str("id", server.ID).Int64("rows_affected", rowsAffected).Msg("Update query executed")

	if rowsAffected == 0 {
		return fmt.Errorf("server with id %s not found for update", server.ID)
	}

	return nil
}

func (s *serverRepositoryImpl) Delete(ctx context.Context, id string) error {
	log := logger.Get()
	log.Debug().Str("id", id).Msg("Executing Delete query")

	res, err := s.db.ExecContext(ctx, config.DeleteServer, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to execute delete query")
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Str("id", id).Int64("rows_affected", rowsAffected).Msg("Delete query executed")

	if rowsAffected == 0 {
		return fmt.Errorf("server with id %s not found for delete", id)
	}

	return nil
}

func (s *serverRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Server, error) {
	log := logger.Get()
	log.Debug().Msg("Executing FindAll query")

	rows, err := s.db.QueryContext(ctx, config.GetListServers)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute find all query")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows in FindAll")
		}
	}(rows)

	var servers []*entity.Server
	for rows.Next() {
		var server entity.Server
		err := rows.Scan(&server.ID, &server.Name, &server.Address, &server.Port, &server.User, &server.Password, &server.CreatedAt, &server.UpdatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan server row in FindAll")
			return nil, err
		}
		servers = append(servers, &server)
	}

	log.Debug().Int("count", len(servers)).Msg("FindAll query finished")
	return servers, nil
}
