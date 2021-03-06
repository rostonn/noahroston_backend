package models

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type User struct {
	ID              int             `json:"-"`
	Email           string          `json:"email"`
	Name            string          `json:"name,omitempty"`
	LastOauth       string          `json:"lastOauth"`
	CreatedTs       string          `json:"createdTs"`
	UpdatedTs       string          `json:"updatedTs"`
	LoginCount      int             `json:"-"`
	UserLoginRecord UserLoginRecord `json:"userLoginRecord"`
}

type UserLoginRecord struct {
	UserID        int     `json:"-"`
	IpAddress     string  `json:"ipAddress"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	OauthProvider string  `json:"oauthProvider"`
	CreatedTs     string  `json:"createdTs,,omitempty"`
}

func (u *User) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO USERS(email, last_oauth,login_count) VALUES('%s','%s', %d)", u.Email, u.LastOauth, 1)
	res, err := db.Exec(statement)
	if err != nil {
		zap.S().Error("INSERT USER ERR", err.Error())
		return err
	}
	// Set user id
	var id int64
	id, err = res.LastInsertId()
	u.ID = int(id)
	return nil
}

func (uL *UserLoginRecord) CreateUserLoginRecord(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO USER_LOGIN_RECORDS("+
		"user_id, ip_address, country_code, country_name, region_code, region_name, city, zip,latitude,"+
		"longitude, oauth_provider) VALUES(%d,'%s','%s','%s','%s','%s','%s', '%s', %f, %f, '%s')",
		uL.UserID, uL.IpAddress, uL.CountryCode, uL.CountryName, uL.RegionCode,
		uL.RegionName, uL.City, uL.Zip, uL.Latitude, uL.Longitude, uL.OauthProvider)

	_, err := db.Exec(statement)
	if err != nil {
		zap.S().Error("INSERT USER LOGIN ERR", err)
		return err
	}
	return nil
}

func (u *User) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE USERS SET last_oauth='%s',login_count=login_count+1 WHERE id=%d", u.LastOauth, u.ID)
	_, err := db.Exec(statement)
	if err != nil {
		zap.S().Error("UPDATE USER ERR", err)
		return err
	}
	return nil
}

func (u *User) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id,email,created_ts,updated_ts,login_count FROM USERS where email='%s'", u.Email)

	row := db.QueryRow(statement)
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedTs, &u.UpdatedTs, &u.LoginCount); err != nil {
		zap.S().Error("User Query Error", err)
		return err
	}

	return nil
}

func (u *User) LoginUser(db *sql.DB) error {
	err := u.getUser(db)
	if err != nil {
		err := u.createUser(db)
		if err != nil {
			return err
		}
	} else {
		err := u.updateUser(db)
		if err != nil {
			return err
		}
	}
	return nil
}
