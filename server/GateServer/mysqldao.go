package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDao struct represents the DAO for MySQL operations
type MysqlDao struct {
	db *sql.DB
}

// NewMysqlDao creates a new instance of MysqlDao
func NewMysqlDao() (*MysqlDao, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", MysqlUser, MysqlPasswd, MysqlHost, MysqlPort, MysqlSchema))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MysqlDao{db: db}, nil
}

// Close closes the database connection
func (dao *MysqlDao) Close() {
	dao.db.Close()
}

// RegUser registers a new user in the database
func (dao *MysqlDao) RegUser(name, email, pwd string) (int, error) {
	tx, err := dao.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// Check if email already exists
	var emailExists int
	err = tx.QueryRow("SELECT 1 FROM user WHERE email = ?", email).Scan(&emailExists)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}
	if emailExists > 0 {
		return 0, nil // Email exists
	}

	// Check if name already exists
	var nameExists int
	err = tx.QueryRow("SELECT 1 FROM user WHERE name = ?", name).Scan(&nameExists)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}
	if nameExists > 0 {
		return 0, nil // Name exists
	}

	// Insert new user and get the new auto-increment ID
	result, err := tx.Exec("INSERT INTO user (name, email, pwd) VALUES (?, ?, ?)", name, email, pwd)
	if err != nil {
		return -1, err
	}

	// Get the new user ID
	newID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return int(newID), nil
}

// UpdatePwd updates user password
func (dao *MysqlDao) UpdatePwd(name, newpwd string) error {
	_, err := dao.db.Exec("UPDATE user SET pwd = ? WHERE name = ?", newpwd, name)
	if err != nil {
		return err
	}
	return nil
}

// CheckPwd checks if the provided password matches the user's password
func (dao *MysqlDao) CheckPwd(email, pwd string) (bool, error) {
	var storedPwd string
	err := dao.db.QueryRow("SELECT pwd FROM user WHERE email = ?", email).Scan(&storedPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}
	return pwd == storedPwd, nil
}

// CheckEmail checks if the provided email matches the user's email
func (dao *MysqlDao) CheckEmail(name, email string) (bool, error) {
	var storedEmail string
	err := dao.db.QueryRow("SELECT email FROM user WHERE name = ?", name).Scan(&storedEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}
	return email == storedEmail, nil
}

// func main() {
// 	dao, err := NewMysqlDao()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize MySQL connection: %v", err)
// 	}
// 	defer dao.Close()

// 	// Example usage
// 	newUserID, err := dao.RegUser("John Doe", "john@example.com", "password123")
// 	if err != nil {
// 		log.Println("Failed to register user:", err)
// 	} else {
// 		fmt.Println("New User ID:", newUserID)
// 	}

// 	err = dao.UpdatePwd("John Doe", "newpassword456")
// 	if err != nil {
// 		log.Println("Failed to update password:", err)
// 	}

// 	validPwd, err := dao.CheckPwd("john@example.com", "newpassword456")
// 	if err != nil {
// 		log.Println("Failed to check password:", err)
// 	} else {
// 		fmt.Println("Password is valid:", validPwd)
// 	}

// 	validEmail, err := dao.CheckEmail("John Doe", "john@example.com")
// 	if err != nil {
// 		log.Println("Failed to check email:", err)
// 	} else {
// 		fmt.Println("Email is valid:", validEmail)
// 	}
// }
