package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sharefood/internal/entity"
	"sharefood/pkg/logger"
	"sharefood/pkg/postgres"
)

type User interface {
	List(context.Context) ([]entity.User, error)
	GetByID(context.Context, int64) (entity.User, error)
	GetByEmail(context.Context, string) (entity.User, error)
	Create(context.Context, *entity.User) error
	IsRegistered(context.Context, string) bool
}

type userImplementation struct {
	conn postgres.Adapter
}

func NewUserRepository(conn postgres.Adapter) User {
	return &userImplementation{conn}
}

// Get all users in the repository
func (r userImplementation) List(ctx context.Context) (users []entity.User, err error) {
	query := "SELECT id_user, name, email, password, image_url from users"
	rows, err := r.conn.QueryRows(ctx, query)
	fmt.Println(err)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.ImageUrl,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Get single user by ID
func (r userImplementation) GetByID(ctx context.Context, id int64) (user entity.User, err error) {
	query := `
		SELECT name
		FROM users
		WHERE id = $1
	`

	row := r.conn.QueryRow(ctx, query, id)
	err = row.Scan(
		&user.Name,
	)
	if err != nil {
		err = fmt.Errorf("scanning user %w", err)
		return entity.User{}, err
	}

	return user, nil
}

// Register User
func (r userImplementation) Create(ctx context.Context, user *entity.User) (err error) {
	query := `
	INSERT INTO users(id_user, email, name, phone_number, password) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.conn.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PhoneNumber,
		user.Password,
	)

	if err != nil {
		fmt.Println(fmt.Errorf("executing query: %w", err))
		return err
	}

	return nil
}

// Get single user by email
func (r userImplementation) GetByEmail(ctx context.Context, email string) (user entity.User, err error) {
	query := `
		SELECT id_user, name, email, phone_number, password, image_url
		FROM users
		WHERE (email = $1) AND (deleted_at IS NULL)
	`

	row := r.conn.QueryRow(ctx, query, email)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.ImageUrl,
	)

	if err != nil {
		err = fmt.Errorf("scanning user %w", err)
		return entity.User{}, err
	}

	return user, nil
}

// Check is user registered by querying 1 matching email
func (r userImplementation) IsRegistered(ctx context.Context, email string) bool {
	query := `
		SELECT name
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	var name string
	err := r.conn.QueryRow(ctx, query, email).Scan(&name)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}
