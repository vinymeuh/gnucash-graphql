// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package gnucash

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadCompressedGnuCashFile(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(file)

	gzfile := "testdata/empty.gnucash.gz"
	defer os.RemoveAll(gzfile)

	cmd := exec.Command("gzip", "-k", "testdata/empty.gnucash")
	cmd.Dir = cwd
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("error while compressing file 'testdata/empty.gnucash': %s", stderr.String())
	}

	_, err := LoadFile(gzfile)
	if err != nil {
		t.Errorf("error reading file 'testdata/empty.gnucash.gz': %s", err)
	}
}

func TestLoadInvalidGnuCashFile(t *testing.T) {
	_, err := LoadFile("gnucashfile_test.go")
	if err == nil {
		t.Error("reading file 'file_test.go' should have failed")
	}
}

func TestLoadDatabaseFromFile(t *testing.T) {
	_, err := LoadFile("i_do_not_exist")
	if err == nil {
		t.Error("reading non existing file should have failed")
	}
}
