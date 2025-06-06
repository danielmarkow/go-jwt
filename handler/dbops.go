package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (ac *AppContext) getUserById(userId int) (*userOut, error) {
	var user userOut
	row := ac.DB.QueryRow("SELECT * FROM users WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ac *AppContext) getUserByEmail(email string) ([]userOut, error) {
	var users []userOut
	rows, err := ac.DB.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	for rows.Next() {
		var user userOut
		err = rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			log.Println("getUserByEmail - error scanning row!")
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (ac *AppContext) createUser(user userIn) error {
	hashedPw, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash pw: %w", err)
	}
	_, err = ac.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, hashedPw)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (ac *AppContext) deleteUserById(userId int) error {
	_, err := ac.DB.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
