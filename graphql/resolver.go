// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package graphql

import (
	"context"

	"github.com/vinymeuh/gnucash-graphql/gnucash"
)

type Resolver struct {
	db *gnucash.Database
}

func (r *queryResolver) RootAccount(ctx context.Context) (*gnucash.Account, error) {
	return r.db.RootAccount, nil
}

func (r *queryResolver) Account(ctx context.Context, id string) (*gnucash.Account, error) {
	return r.db.Account(id), nil
}

func (r *queryResolver) Accounts(ctx context.Context, name string) ([]*gnucash.Account, error) {
	root := r.db.RootAccount
	return root.Descendants(gnucash.AccountFilter{Name: name}), nil
}

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
