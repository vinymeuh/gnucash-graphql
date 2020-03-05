// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

import "testing"

func TestAccount(t *testing.T) {
	acts := Account{
		ID:   "0",
		Name: "Dummy Account",
		Type: "ROOT",
		Children: []*Account{
			{
				ID:   "1",
				Name: "Account 1",
				Type: "BANK",
			},
			{
				ID:   "2",
				Name: "Account 2",
				Type: "BANK",
			},
		},
	}

	root := acts.Descendants(AccountFilter{Type: "ROOT"})
	if len(root) != 1 {
		t.Errorf("too many accounts with Type 'ROOT', expected=1, got=%d", len(root))
	} else {
		if root[0] != &acts {
			t.Error("the root account retrieved by Type is not the expected one")
		}
	}

	root2 := acts.Descendants(AccountFilter{ID: "0", Name: "Dummy Account", Type: "ROOT"})
	if len(root2) != 1 {
		t.Errorf("too many accounts with ID '0', Name 'Dummy Account' and Type 'ROOT', expected=1, got=%d", len(root2))
	} else {
		if root2[0] != &acts {
			t.Error("the root account retrieved by ID x Name x Type is not the expected one")
		}
	}

	root3 := acts.Descendants(AccountFilter{Name: "Dummy Account", Type: "ROOT"})
	if len(root3) != 1 {
		t.Errorf("too many accounts with Name 'Dummy Account' and Type 'ROOT', expected=1, got=%d", len(root2))
	} else {
		if root3[0] != &acts {
			t.Error("the root account retrieved by Name x Type is not the expected one")
		}
	}

	act1 := acts.Descendants(AccountFilter{ID: "1"})
	if len(act1) != 1 {
		t.Errorf("too many accounts with ID 1, expected=1, got=%d", len(act1))
	} else {
		if act1[0] != acts.Children[0] {
			t.Error("the account 1 retrieved by ID is not the expected one")
		}
	}

	act1b := acts.Descendants(AccountFilter{ID: "1", Name: "Account 1"})
	if len(act1b) != 1 {
		t.Errorf("too many accounts with ID 1 and Name 'Account 1', expected=1, got=%d", len(act1b))
	} else {
		if act1b[0] != acts.Children[0] {
			t.Error("the account 1 retrieved by ID x Name is not the expected one")
		}
	}

	act2 := acts.Descendants(AccountFilter{Name: "Account 2"})
	if len(act2) != 1 {
		t.Errorf("too many accounts with Name 'Account 2', expected=1, got=%d", len(act2))
	} else {
		if act2[0] != acts.Children[1] {
			t.Error("the ID 2 account retrieved by ID is not the expected one")
		}
	}

	act2b := acts.Descendants(AccountFilter{ID: "2", Type: "BANK"})
	if len(act2b) != 1 {
		t.Errorf("too many accounts with Name 'Account 2' and Type 'BANK', expected=1, got=%d", len(act2))
	} else {
		if act2[0] != acts.Children[1] {
			t.Error("the ID 2 account retrieved by ID x Type is not the expected one")
		}
	}

	all := acts.Descendants(AccountFilter{})
	if len(all) != 3 {
		t.Errorf("wrong number of accounts retrieved, expected=%d, got=%d", 3, len(all))
	}

	dummy := acts.Descendants(AccountFilter{ID: "DuMmY"})
	if len(dummy) != 0 {
		t.Error("we retrieved a non existing account!")
	}

}
