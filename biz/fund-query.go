package biz

import (
	"database/sql"
	"neolong.me/fundinfo/dbs"
)

func HasFundCodeData(code string) (bool, error) {
	dbConn, err := dbs.OpenDefaultDB()
	if nil != err {
		return false, err
	}

	rows, err := dbConn.Query("select count(*) from fund_worth where fund_code=?", code);
	if nil != err && err != sql.ErrNoRows {
		return false, err
	}
	return nil != err && rows.Next(), nil
}