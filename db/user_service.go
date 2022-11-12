package db

import (
	"context"
	"database/sql"
)

type UserService struct {
	dbConn *sql.DB
}

func NewUserService(dbConn *sql.DB) *UserService {
	return &UserService{dbConn: dbConn}
}

func (s UserService) Create(ctx context.Context, newUser UserCreateDTO) (User, error) {
	var createdUser User

	dbTx, err := s.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer dbTx.Rollback()

	query := `
		INSERT INTO "user"(name) 
		VALUES ($1)
		RETURNING id, name
	`
	err = dbTx.QueryRow(query, newUser.Name).Scan(&createdUser.Id, &createdUser.Name)
	if err != nil {
		return User{}, err
	}

	query = `
		INSERT INTO balance(user_id, amount)
		VALUES ($1,$2)
	`
	_, err = dbTx.Exec(query, createdUser.Id, 1000)
	if err != nil {
		return User{}, err
	}

	if err = dbTx.Commit(); err != nil {
		return User{}, err
	}

	return createdUser, nil
}

func (s UserService) GetAll() ([]User, error) {
	var (
		users []User
		user  User
	)
	query := `SELECT id, name from "user"`
	rows, err := s.dbConn.Query(query)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s UserService) GetUserWithBalance(userId int) (UserWithBalanceOutDTO, error) {
	var (
		userWithBalance UserWithBalanceOutDTO
	)

	query := `
		SELECT u.id, u.name, b.id, b.amount FROM "user" u
		JOIN balance b ON u.id = b.user_id
		WHERE u.id = $1
	`
	err := s.dbConn.QueryRow(query, userId).
		Scan(
			&userWithBalance.Id,
			&userWithBalance.Name,
			&userWithBalance.Balance.Id,
			&userWithBalance.Balance.Amount,
		)
	if err != nil {
		return UserWithBalanceOutDTO{}, err
	}

	return userWithBalance, nil
}
