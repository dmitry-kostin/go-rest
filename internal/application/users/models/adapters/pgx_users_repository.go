package adapters

import (
	"context"
	"github.com/dmitry-kostin/go-rest/internal/application/users/models"
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

func (r PgxUsersRepository) CreateUser(data *models.CreateUserReq) (*models.User, error) {
	ctx := context.Background()
	sqlStatement := `INSERT INTO users 
	(first_name, last_name, email, role) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, email, first_name, last_name, created_at, role`
	row := r.pool.QueryRow(ctx, sqlStatement, data.FirstName, data.LastName, data.Email, data.Role)
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.Role)
	return &user, err
}

func (r PgxUsersRepository) ListUsers() (*[]models.User, error) {
	ctx := context.Background()
	sqlStatement := `select id, email, first_name, last_name, created_at, role from users`
	rows, err := r.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, err
}

func (r PgxUsersRepository) RemoveUser(userId int64) error {
	ctx := context.Background()
	sqlStatement := `delete from users where id = $1`
	_, err := r.pool.Exec(ctx, sqlStatement, userId)
	return err
}

func (r PgxUsersRepository) GetUser(userId int64) (*models.User, error) {
	ctx := context.Background()
	sqlStatement := `select id, email, first_name, last_name, created_at, role from users where id = $1`
	row := r.pool.QueryRow(ctx, sqlStatement, userId)
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.Role)
	return &user, err
}
