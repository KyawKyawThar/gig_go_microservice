package sqldb

import "auth_service/db/mdata"

func (sql *SqlDB) GetCurrentUserById(authId string) (*mdata.Auth, error) {

	var auth mdata.Auth

	cols := []string{
		"id",
		"username",
	}
	err := sql.FetchOne(&auth, "auths", cols, "id=?", authId)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}
