// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

type Database struct {
	RootAccount *Account
	Statistics  DBStatistics
	Warnings    []error
	actsIndex   map[string]*Account
}

type DBStatistics struct {
	File         string
	Accounts     int
	Transactions int
	LoadDate     string
	LoadTime     string
}

func (db *Database) Account(id string) *Account {
	return db.actsIndex[id]
}
