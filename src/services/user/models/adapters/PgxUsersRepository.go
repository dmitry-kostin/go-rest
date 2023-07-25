package adapters

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/dmitry-kostin/go-rest/src/services/user/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxUsersRepository struct {
	pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) models.UserRepository {
	return &PgxUsersRepository{pool: pool}
}

func (s *PgxUsersRepository) CreateUser(user *models.User) error {
	wrapWith := "[PgxUsersRepository.CreateUser]"
	ctx := context.Background()
	sqlStatement := `INSERT INTO users 
	(id, identity_id, email, role, first_name, last_name, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, sqlStatement, user.Id, user.IdentityId, user.Email, user.Role, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return s.handleError(err, wrapWith)
	}
	return nil
}

func (s *PgxUsersRepository) ListUsers() ([]*models.User, error) {
	wrapWith := "[PgxUsersRepository.ListUsers]"
	ctx := context.Background()
	sqlStatement := `select id, identity_id, email, role, first_name, last_name, created_at, updated_at from users`
	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, s.handleError(err, wrapWith)
	}
	users := make([]*models.User, 0)
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.IdentityId, &user.Email, &user.Role, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, s.handleError(err, wrapWith)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *PgxUsersRepository) RemoveUser(userId models.EntityId) error {
	wrapWith := "[PgxUsersRepository.RemoveUser]"
	ctx := context.Background()
	sqlStatement := `delete from users where id = $1`
	_, err := s.pool.Exec(ctx, sqlStatement, userId)
	if err != nil {
		return s.handleError(err, wrapWith)
	}
	return nil
}

func (s *PgxUsersRepository) GetUser(userId models.EntityId) (*models.User, error) {
	wrapWith := "[PgxUsersRepository.GetUser]"
	ctx := context.Background()
	sqlStatement := `select id, identity_id, email, role, first_name, last_name, created_at, updated_at from users where id = $1`
	row := s.pool.QueryRow(ctx, sqlStatement, userId)
	var user models.User
	err := row.Scan(&user.Id, &user.IdentityId, &user.Email, &user.Role, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, s.handleError(err, wrapWith)
	}
	return &user, nil
}

func (s *PgxUsersRepository) handleError(err error, wrapWith string) error {
	var duplicateEntryError = &pgconn.PgError{Code: "23505"}
	if errors.As(err, &duplicateEntryError) {
		return pkg.AnnotateErrorWithDetail(errors.New("duplicate email address"), pkg.ErrDuplicate, wrapWith, "Provided email address is already in use")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return pkg.AnnotateErrorWithDetail(err, pkg.ErrNotFound, wrapWith, "User not found")
	}
	return pkg.AnnotateError(err, pkg.ErrDatabaseError, wrapWith)
}
