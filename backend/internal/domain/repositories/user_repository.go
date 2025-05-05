package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Delete(ctx context.Context, id int) error
	FindByNameLike(ctx context.Context, name string) ([]*models.User, error)
	UpdateXpAndLevel(ctx context.Context, user *models.User) error
}

type UserRepository struct {
	db *pgx.Conn
}

// // xpRequiredForLevel — формула расчета XP для уровня
// func xpRequiredForLevel(level int) int {
// 	return 100 * level * level
// }

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, name, surname, email, password, role, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRow(ctx, query,
		user.Username, user.Name, user.Surname, user.Email, user.Password, user.Role).
		Scan(&user.ID)

	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, name, surname, email, password, role, level, xp, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.Level, &user.Xp, &user.CreatedAt, &user.UpdatedAt)
	fmt.Println(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, name, surname, email, password, role, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, name, surname, email, password, role, created_at, updated_at FROM users WHERE username = $1`

	err := r.db.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) FindByNameLike(ctx context.Context, usernameLike string) ([]*models.User, error) {
	query := `SELECT id, username, name, surname, email, password, role, created_at, updated_at FROM users WHERE username ILIKE $1`
	rows, err := r.db.Query(ctx, query, "%"+usernameLike+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) UpdateXpAndLevel(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET xp = $1, level = $2, updated_at = NOW()
		WHERE id = $3;
	`
	
	_, err := r.db.Exec(ctx, query, user.Xp, user.Level, user.ID)
	return err
}