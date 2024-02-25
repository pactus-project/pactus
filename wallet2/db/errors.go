package db

import "errors"

var (
	ErrCouldNotOpenDatabase          = errors.New("could not open database")
	ErrCouldNotCreateTable           = errors.New("could not create table")
	ErrCouldNotInsertRecordIntoTable = errors.New("could not insert record into table")
	ErrCouldNotUpdateRecordIntoTable = errors.New("could not update record into table")
	ErrCouldNotFindRecord            = errors.New("could not find record")
	ErrCouldNotFindTotalRecords      = errors.New("could not find total records")
)
