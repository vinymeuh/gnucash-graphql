// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

// Account is a node of the accounts hierarchy.
// Each account has its own list of transactions
type Account struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Parent       *Account       `json:"-"`
	Children     []*Account     `json:"-"`
	Transactions []*Transaction `json:"-"`
}

// Transaction keeps data for a transaction
type Transaction struct {
	Num   string  `json:"-"`
	Date  string  `json:"-"` // YYYY-MM-DD
	Value float64 `json:"-"`
}

// WalkAccountFunc is the type of the function called for each account visited by WalkBFS
type WalkAccountFunc func(act *Account) bool

// WalkBFS traverses the tree of accounts using Breadth-first search algorithm.
// Starting at node a, returns the list of accounts for which walkFunc is true.
func (a *Account) WalkBFS(walkFunc WalkAccountFunc) []*Account {
	acts := make([]*Account, 0)

	queue := make([]*Account, 0)
	queue = append(queue, a)
	for len(queue) > 0 {
		act := queue[0]
		if walkFunc(act) == true {
			acts = append(acts, act)
		}
		queue = queue[1:]
		if len(act.Children) > 0 {
			for _, child := range act.Children {
				queue = append(queue, child)
			}
		}
	}

	return acts
}

// AccountFilter is the type to express filter condition for Descendants method.
// Note that it is not possible to filter accounts with Name empty
type AccountFilter struct {
	ID   string
	Name string
	Type string
}

// Descendants returns a list of sub-accounts matching the filter,
// using account a as the starting point of the search
func (a *Account) Descendants(filter AccountFilter) []*Account {
	var filterFunc WalkAccountFunc

	switch {
	case filter.ID != "" && filter.Name != "" && filter.Type != "":
		filterFunc = func(act *Account) bool {
			return act.ID == filter.ID && act.Name == filter.Name && act.Type == filter.Type
		}
	case filter.ID != "" && filter.Name == "" && filter.Type != "":
		filterFunc = func(act *Account) bool {
			return act.ID == filter.ID && act.Type == filter.Type
		}
	case filter.ID != "" && filter.Name != "" && filter.Type == "":
		filterFunc = func(act *Account) bool {
			return act.ID == filter.ID && act.Name == filter.Name
		}
	case filter.ID != "" && filter.Name == "" && filter.Type == "":
		filterFunc = func(act *Account) bool {
			return act.ID == filter.ID
		}
	case filter.ID == "" && filter.Name != "" && filter.Type != "":
		filterFunc = func(act *Account) bool {
			return act.Name == filter.Name && act.Type == filter.Type
		}
	case filter.ID == "" && filter.Name == "" && filter.Type != "":
		filterFunc = func(act *Account) bool {
			return act.Type == filter.Type
		}
	case filter.ID == "" && filter.Name != "" && filter.Type == "":
		filterFunc = func(act *Account) bool {
			return act.Name == filter.Name
		}
	case filter.ID == "" && filter.Name == "" && filter.Type == "":
		filterFunc = func(act *Account) bool {
			return true
		}
	}

	return a.WalkBFS(filterFunc)
}
