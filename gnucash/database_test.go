// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

import "testing"

func TestLoadDatabase(t *testing.T) {
	db, err := LoadFile("testdata/empty.gnucash")
	if err != nil {
		t.Fatalf("error reading file 'testdata/empty.gnucash': %s", err)
	}

	if len(db.Warnings) > 0 {
		t.Errorf("db loaded with warnings: %+v", db.Warnings)
	}

	root := db.Account(db.RootAccount.ID)
	if root.Type != "ROOT" {
		t.Error("problem while searching the root account by ID")
	}

	dummy := db.Account("DuMmY")
	if dummy != nil {
		t.Error("problem while searching a non existing account by ID")
	}
}
