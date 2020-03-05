// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

import (
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type xmlCountData struct {
	Key   string `xml:"type,attr"`
	Value int    `xml:",chardata"`
}

type xmlAccount struct {
	Name     string `xml:"name"`
	ID       string `xml:"id"`
	Type     string `xml:"type"`
	ParentID string `xml:"parent"`
	Parent   *xmlAccount
	Children []*xmlAccount
}

type xmlTransaction struct {
	Num        string     `xml:"num"`
	DatePosted string     `xml:"date-posted>date"`
	Splits     []xmlSplit `xml:"splits>split"`
}

type xmlSplit struct {
	Value     string `xml:"value"`
	AccountID string `xml:"account"`
}

// LoadFile loads a Database from a GnuCash file (compressed or not)
func LoadFile(path string) (*Database, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var db *Database
	zr, err := gzip.NewReader(f)
	switch err {
	case nil:
		defer zr.Close()
		db, err = Load(zr)
	case gzip.ErrHeader: // uncompressed file
		f.Seek(0, 0)
		db, err = Load(f)
	default:
		return nil, err
	}

	db.Statistics.File = path
	return db, err
}

// Load loads a Database from a Gnucash XML document stream
func Load(r io.Reader) (*Database, error) {
	db := &Database{}

	var (
		acntReaded int
		trnsReaded int
	)

	t1 := time.Now()
	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Space != "http://www.gnucash.org/XML/gnc" {
				continue
			}

			//---- COUNT-DATA ----//
			if se.Name.Local == "count-data" {
				var cd xmlCountData
				decoder.DecodeElement(&cd, &se)

				switch cd.Key {
				case "account":
					db.Statistics.Accounts = cd.Value
					db.actsIndex = make(map[string]*Account, cd.Value)
				case "transaction":
					db.Statistics.Transactions = cd.Value
				}
				continue
			}

			//---- ACCOUNT ----//
			if se.Name.Local == "account" {
				var xmlact xmlAccount
				decoder.DecodeElement(&xmlact, &se)
				acntReaded++

				// I hope Root Account is always the first account encountered
				if db.RootAccount == nil {
					if xmlact.Type == "ROOT" && xmlact.Name == "Root Account" {
						db.RootAccount = &Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type}
						db.actsIndex[xmlact.ID] = db.RootAccount
						continue
					}
					return nil, errors.New("Unable to initialize accounts hierarchy with Root Account")
				}

				// Attach this node to the accounts tree
				parent := db.actsIndex[xmlact.ParentID]
				if parent == nil {
					db.Warnings = append(db.Warnings, fmt.Errorf("ParentID not found in index for Account '%s'", xmlact.Name))
					continue
				}
				act := Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type, Parent: parent}
				parent.Children = append(parent.Children, &act)

				db.actsIndex[xmlact.ID] = &act
				continue
			}

			//---- TRANSACTION ----//
			if se.Name.Local == "transaction" {
				var xtrn xmlTransaction
				decoder.DecodeElement(&xtrn, &se)
				trnsReaded++
				for _, split := range xtrn.Splits {
					act := db.actsIndex[split.AccountID]
					if act == nil {
						db.Warnings = append(db.Warnings, fmt.Errorf("Account '%s' not found in index for transaction", split.AccountID))
						continue
					}
					trn := Transaction{
						Num:   xtrn.Num,
						Date:  strings.TrimSpace(strings.Split(xtrn.DatePosted, " ")[0]), // '2014-07-30 00:00:00 +0200', we keep only '2014-07-30'
						Value: stringToFloat(split.Value),
					}
					act.Transactions = append(act.Transactions, &trn)
				}
			}

			// Skip all accounts and transactions templates used in schedule action
			if se.Name.Local == "template-transactions" {
				decoder.Skip()
			}
		}
	}
	t2 := time.Now()
	db.Statistics.LoadTime = fmt.Sprintf("%s", t2.Sub(t1))

	// Compare readed vs expected
	if acntReaded != db.Statistics.Accounts {
		db.Warnings = append(db.Warnings, fmt.Errorf("Read %d accounts when %d were expected", acntReaded, db.Statistics.Accounts))
	}

	if trnsReaded != db.Statistics.Transactions {
		db.Warnings = append(db.Warnings, fmt.Errorf("Read %d transactions when %d were expected", trnsReaded, db.Statistics.Transactions))
	}

	db.Statistics.LoadDate = time.Now().Format(time.RFC3339)
	return db, nil
}

func stringToFloat(v string) float64 {
	i := strings.Split(v, "/")
	n, _ := strconv.ParseFloat(i[0], 10)
	d, _ := strconv.ParseFloat(i[1], 10)
	return n / d
}
