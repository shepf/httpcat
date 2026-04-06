package common

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin"
)

func HashPassword(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password cannot be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func IsBcryptHash(hash string) bool {
	return strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$")
}

func verifyLegacyPassword(password, salt, hash string) (bool, error) {
	t := sha1.New()
	if _, err := io.WriteString(t, password+salt); err != nil {
		return false, err
	}
	return fmt.Sprintf("%x", t.Sum(nil)) == hash, nil
}

func VerifyPassword(user *User, password string) (valid bool, legacy bool, err error) {
	if user == nil {
		return false, false, errors.New("user is nil")
	}

	if IsBcryptHash(user.Password) {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return false, false, nil
			}
			return false, false, err
		}
		return true, false, nil
	}

	valid, err = verifyLegacyPassword(password, user.Salt, user.Password)
	return valid, true, err
}

func MustChangePassword(user *User) bool {
	return false
}

func UpdateUserPasswordCache(username, hashedPassword, salt string, updateTime int64) {
	UserLock.Lock()
	defer UserLock.Unlock()

	user, ok := UserTable[username]
	if !ok || user == nil {
		return
	}

	user.Password = hashedPassword
	user.Salt = salt
	user.PasswordUpdateTime = updateTime
}

func UpgradeLegacyPasswordHash(user *User, plainPassword string) error {
	if user == nil {
		return errors.New("user is nil")
	}

	hashedPassword, err := HashPassword(plainPassword)
	if err != nil {
		return err
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	if err := db.Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"password": hashedPassword,
		"salt":     "",
	}).Error; err != nil {
		return err
	}

	UpdateUserPasswordCache(user.Username, hashedPassword, "", user.PasswordUpdateTime)
	user.Password = hashedPassword
	user.Salt = ""
	return nil
}
