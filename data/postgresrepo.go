package data

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type PostgressRepository struct {
	db     *sqlx.DB
	logger hclog.Logger
}

func NewPostgressRepository(db *sqlx.DB, logger hclog.Logger) *PostgressRepository {
	return &PostgressRepository{db, logger}
}

// Create inserts the given user into the database
func (repo *PostgressRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.logger.Info("Create user", hclog.Fmt("%#v", user))
	query := "INSERT INTO USERS (id, email, username, password, tokenhash, createdat, updatedat) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err := repo.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Password, user.TokenHash, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves the user object having the given email, else returns error
func (repo *PostgressRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	repo.logger.Debug("querying for user with email", email)
	query := "SELECT * FROM USERS WHERE email = $1"
	var user User
	if err := repo.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}

	repo.logger.Debug("read users", hclog.Fmt("%#v", user))

	return &user, nil
}

// GetUserByID retrieves the user object having the given ID, else returns error
func (repo *PostgressRepository) GetUserById(ctx context.Context, id string) (*User, error) {
	repo.logger.Debug("querying for user with id", id)
	query := "SELECT * FROM users WHERE id = $1"
	var user User
	if err := repo.db.GetContext(ctx, query, id); err != nil {
		return nil, err
	}

	repo.logger.Debug("read users", hclog.Fmt("%#v", user))

	return &user, nil
}

// // UpdateUsername updates the username of the given user
func (repo *PostgressRepository) UpdateUsername(ctx context.Context, user *User, username string) error {
	user.UpdatedAt = time.Now()
	query := "UPDATE users SET username = $1, updatedat = $2 WHERE id = $3"
	_, err := repo.db.ExecContext(ctx, query, user.Username, user.UpdatedAt, user.ID)
	return err
}

// UpdateUserVerificationStatus updates user verification status to true
func (repo *PostgressRepository) UpdateUserVerificationStatus(ctx context.Context, email string, status bool) error {
	query := "UPDATE users SET status = $1 WHERE email = $2"
	_, err := repo.db.ExecContext(ctx, query, status, email)
	return err
}

// StoreMailVerificationData adds a mail verification data to db
func (repo *PostgressRepository) StoreMailVerificationData(ctx context.Context, verificationData *VerificationData) error {
	query := "INSERT INTO verifications(email, code, expiresat, type) VALUES($1, $2, $3, $4)"
	_, err := repo.db.ExecContext(ctx, query, verificationData.Email, verificationData.Code, verificationData.ExpiresAt, verificationData.Type)
	return err
}

// GetMailVerificationCode retrieves the stored verification code.
func (repo *PostgressRepository) GetMailVerificationCode(ctx context.Context, email string, verificationDataType VerificationDataType) (error, *VerificationData) {
	query := "SELECT * FROM verifications WHERE email = $1 AND type = $2"
	var verificationData VerificationData
	if err := repo.db.GetContext(ctx, &verificationData, query, email, verificationDataType); err != nil {
		return err, nil
	}

	return nil, &verificationData
}

// DeleteMailVerificationData deletes a used verification data
func (repo *PostgressRepository) DeleteMailVerificationData(ctx context.Context, email string, verificationDataType VerificationDataType) error {
	query := "DELETE FROM verifications WHERE email = $1 and type = $2"
	_, err := repo.db.ExecContext(ctx, query, email, verificationDataType)
	return err
}

// UpdatePassword updates the user password
func (repo *PostgressRepository) UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error {
	query := "UPDATE users SET password = $1, tokenhash = $2 WHERE id = $3"
	_, err := repo.db.ExecContext(ctx, query, password, tokenHash, userID)
	return err
}
