package adapters

import (
	"context"
	"github.com/dmitry-kostin/go-rest/src/services/users/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PgxUsersRepository struct {
	pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) models.UserRepository {
	if pool == nil {
		log.Fatal("Missing pool connection")
	}
	return &PgxUsersRepository{pool: pool}
}

func (r PgxUsersRepository) CreateUser(user *models.User) error {
	ctx := context.Background()
	sqlStatement := `INSERT INTO users 
	(id, identity_id, email, role, first_name, last_name, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, sqlStatement, user.Id, user.IdentityId, user.Email, user.Role, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (r PgxUsersRepository) ListUsers() (*[]models.User, error) {
	ctx := context.Background()
	sqlStatement := `select id, identity_id, email, role, first_name, last_name, created_at, updated_at from users`
	rows, err := r.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.IdentityId, &user.Email, &user.Role, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, err
}

func (r PgxUsersRepository) RemoveUser(userId models.EntityId) error {
	ctx := context.Background()
	sqlStatement := `delete from users where id = $1`
	_, err := r.pool.Exec(ctx, sqlStatement, userId)
	return err
}

func (r PgxUsersRepository) GetUser(userId models.EntityId) (*models.User, error) {
	ctx := context.Background()
	sqlStatement := `select id, identity_id, email, role, first_name, last_name, created_at, updated_at from users where id = $1`
	row := r.pool.QueryRow(ctx, sqlStatement, userId)
	var user models.User
	err := row.Scan(&user.Id, &user.IdentityId, &user.Email, &user.Role, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}
