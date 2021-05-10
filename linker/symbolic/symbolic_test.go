package symbolic

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLinkFromDir(t *testing.T) {
	log := logrus.New()
	testCases := []struct {
		desc                            string
		in                              string
		out                             string
		recursive                       bool
		whitelisted                     []string
		numUnlinkedAttempts, totalFiles int
		err                             error
	}{
		{
			desc:                "dotfiles are linked and theres no errors",
			in:                  "./fixtures/",
			out:                 createTmpDir(t, "./"),
			whitelisted:         []string{"Whitelisted"},
			recursive:           true,
			numUnlinkedAttempts: 1,
			totalFiles:          5,
			err:                 nil,
		},
		{
			desc:                "git files are ignore even whitelisted",
			in:                  "./fixtures/",
			out:                 createTmpDir(t, "./"),
			whitelisted:         []string{"Whitelisted", ".gitignore"},
			recursive:           true,
			numUnlinkedAttempts: 1,
			totalFiles:          5,
			err:                 nil,
		},
		// {
		// 	desc:        "duplicate filename whitelisted ignored - linked success",
		// 	in:          mustAbsPath(t, "./fixtures/dotfiles/"),
		// 	out:         mustAbsPath(t, "./fixtures/symlinked/"),
		// 	whitelisted: []string{"Whitelisted", "Whitelisted"},
		// 	want:        0,
		// 	err:         nil,
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			a := assert.New(t)

			op, err := NewDirOp(log, tC.in, tC.out, tC.recursive, tC.whitelisted...)
			a.NoError(err)
			unlinked, total, err := op.Link()
			a.ErrorIs(err, tC.err)
			a.Equal(tC.numUnlinkedAttempts, unlinked)
			a.Equal(tC.totalFiles, total)
			time.Sleep(2 * time.Second)
			t.Cleanup(removeFiles(t, tC.out))
		})
	}
}

func removeFiles(t *testing.T, dir string) func() {
	t.Helper()
	return func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}
}

func createTmpDir(t *testing.T, str string) (dname string) {
	t.Helper()
	dname, err := ioutil.TempDir(str, "*")
	if err != nil {
		t.Fatal(err)
	}
	return dname
}

func mustAbsPath(t *testing.T, in string) string {
	t.Helper()
	f, err := filepath.Abs(in)
	if err != nil {
		t.Error(err)
	}
	return f
}
