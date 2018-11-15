package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID               int               `json:"-"`
	Email            string            `json:"email"`
	Name             string            `json:"name"`
	LastOauth        string            `json:"lastOauth"`
	CreatedTs        string            `json:"createdTs"`
	UpdatedTs        string            `json:"updatedTs"`
	LoginCount       int               `json:"-"`
	UserLoginRecords []UserLoginRecord `json:"userLoginRecord"`
}

type UserLoginRecord struct {
	UserID        int     `json:"-"`
	IpAddress     []byte  `json:"ipAddress"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	OauthProvider string  `json:"oauthProvider"`
	CreatedTs     string  `json:"createdTs`
}

func (u *User) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO USERS(email, last_oauth,login_count) VALUES('%s','%s', %d)", u.Email, u.LastOauth, 1)
	res, err := db.Exec(statement)
	if err != nil {
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

	fmt.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE USERS SET last_oauth='%s',login_count=login_count+1 WHERE id=%d", u.LastOauth, u.ID)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id,email,last_oauth,created_ts,updated_ts,login_count FROM USERS where email='%s'", u.Email)
	fmt.Println("usr", statement)

	row := db.QueryRow(statement)
	if err := row.Scan(&u.ID, &u.Email, &u.LastOauth, &u.CreatedTs, &u.UpdatedTs, &u.LoginCount); err != nil {
		fmt.Println("User Query Error", err)
		return err
	}

	return nil
}

// func getUserLoginRecords(db *sql.DB) ([]UserLoginRecord, error) {

// }

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
