package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Fullname  string   `json:"fullname"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Role      Role     `json:"role"`
	ImgURL    ImgURL   `json:"image_url"`
}

type Password struct {
	Text *string
	Hash []byte
}

type ImgURL struct {
	UserID    int64  `json:"user_id"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash

	return nil
}

func (p *Password) Compared(password string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(password))
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, username, fullname, email, password, is_active, users.created_at, users.updated_at, role, 
		COALESCE(img.user_id,0) AS user_id,
		COALESCE(img.image_url,'') AS image_url,
		COALESCE(img.created_at,NOW()) AS created_at,
		COALESCE(img.updated_at,NOW()) AS updated_at,
		r.level
		FROM users
		LEFT JOIN image_profile img ON (users.id = img.user_id)
		LEFT JOIN roles r ON (users.role = r.name)
		WHERE users.id = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	user := new(User)

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID).Scan(
		&user.ID,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.Password.Hash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.Name,
		&user.ImgURL.UserID,
		&user.ImgURL.ImageURL,
		&user.ImgURL.CreatedAt,
		&user.ImgURL.UpdatedAt,
		&user.Role.Level,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStorage) insertUser(ctx context.Context, tx *sql.Tx, user *User) (*User, error) {
	query := `
		INSERT INTO users (username, fullname, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	if user.Role.Name == "" {
		user.Role.Name = "user"
	}
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Fullname,
		user.Email,
		user.Password.Hash,
		user.Role.Name,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return nil, ErrDuplicateUsername
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStorage) insertImage(ctx context.Context, tx *sql.Tx, userID int64, imgURL ImgURL) error {
	query := `
		INSERT INTO image_profile (user_id, image_url)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, userID, imgURL.ImageURL)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStorage) CreateUser(ctx context.Context, u *User, img *ImgURL) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.insertUser(ctx, tx, u)
		if err != nil {
			return err
		}

		if err := s.insertImage(ctx, tx, user.ID, *img); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT users.id, username, fullname, email, password, is_active, users.created_at, users.updated_at, role, 
		COALESCE(img.user_id,0) AS user_id,
		COALESCE(img.image_url,'') AS image_url,
		COALESCE(img.created_at,NOW()) AS created_at,
		COALESCE(img.updated_at,NOW()) AS updated_at,
		r.level
		FROM users
		LEFT JOIN image_profile img ON (users.id = img.user_id)
		LEFT JOIN roles r ON (users.role = r.name)
		WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	user := new(User)

	err := s.db.QueryRowContext(
		ctx,
		query,
		email).Scan(
		&user.ID,
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.Password.Hash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.Name,
		&user.ImgURL.UserID,
		&user.ImgURL.ImageURL,
		&user.ImgURL.CreatedAt,
		&user.ImgURL.UpdatedAt,
		&user.Role.Level,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
