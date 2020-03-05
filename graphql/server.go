// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/vinymeuh/gnucash-graphql/gnucash"
)

// NewHandler returns a new graphql endpoint handler.
func NewHandler(db *gnucash.Database) http.Handler {
	return handler.GraphQL(NewExecutableSchema(Config{
		Resolvers: &Resolver{
			db: db,
		},
	}))
}
