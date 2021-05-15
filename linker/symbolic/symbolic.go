package symbolic

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/blixenkrone/dotgo/linker"
)

var (
	// ErrInvalidDotfile is returned if the file is not a dotfile i.e. it's a dir or non dot prefixed
	ErrInvalidDotfile = errors.New("not a valid dotfile")

	blacklisted = map[string]struct{}{
		".git":       {},
		".gitignore": {},
		".DS_Store":  {},
	}
)

type DirOperation struct {
	// file link and linked duplicate - in would normally be in the usr $HOME dir
	log                 linker.Logger
	inputDir, outputDir string
	whitelisted         map[string]bool
	recursive           bool
}

func NewDirOp(log linker.Logger, inDir, outDir string, recursive bool, whitelisted ...string) (*DirOperation, error) {
	absDirSrc, err := filepath.Abs(inDir)
	if err != nil {
		return nil, err
	}
	absDirDest, err := filepath.Abs(outDir)
	if err != nil {
		return nil, err
	}

	var wl = make(map[string]bool)
	for _, file := range whitelisted {
		if _, ok := wl[file]; ok {
			// file was targeted twice for some reason
			continue
		}
		wl[file] = true
	}
	return &DirOperation{log, absDirSrc, absDirDest, wl, recursive}, nil
}

// Attempts to create a symbolic link from the specified operation src and dest paths.
func (o *DirOperation) Link() (int, int, error) {
	return o.linkFromDir()
}

// Returns number of files not linked, total files found and err
func (o *DirOperation) linkFromDir() (int, int, error) {
	dirEntries, err := os.ReadDir(o.inputDir)
	totalFiles := len(dirEntries)
	unLinkedFiles := len(dirEntries)
	if err != nil {
		return unLinkedFiles, totalFiles, err
	}

	for _, entry := range dirEntries {
		if _, ok := blacklisted[entry.Name()]; ok {
			o.log.Warnf("file '%s' is blacklisted", entry.Name())
			continue
		}
		d := dotFile{
			entry:       entry,
			whitelisted: false,
		}
		if v, ok := o.whitelisted[entry.Name()]; ok {
			d.whitelisted = v
		}

		if !o.recursive {
			if entry.IsDir() {
				o.log.Infof("entry is a dir %s - skipping...", entry.Name())
				continue
			}
		}
		if err := d.linkFileSymbolic(o.inputDir, o.outputDir); err != nil {
			if errors.Is(err, ErrInvalidDotfile) {
				o.log.Warnf("file %s is not valid dotfile", d.Name())
				continue
			}
			if errors.Is(err, os.ErrExist) {
				o.log.Warnf("file already exists %s", d.Name())
				continue
			}
			o.log.Errorf("file '%s' link error: %v", d.Name(), err)
			return unLinkedFiles, len(dirEntries), err
		}
		unLinkedFiles--
	}
	return unLinkedFiles, len(dirEntries), nil
}

type dotFile struct {
	// name        string
	entry       fs.DirEntry
	whitelisted bool
}

func (d dotFile) linkFileSymbolic(srcDir, destDir string) error {
	srcPath := filepath.Join(srcDir, d.Name())
	if d.hasDestFile(srcPath) {
		// should append file bytes to file in future release?
		return &SymbolicLinkError{err: os.ErrExist, fname: d.Name()}
	}

	if !d.hasDotPrefix() && !d.whiteListed() {
		return &SymbolicLinkError{err: ErrInvalidDotfile, fname: d.Name()}
	}

	destPath := filepath.Join(destDir, d.Name())

	if err := os.Symlink(srcPath, destPath); err != nil {
		return &SymbolicLinkError{err: err, fname: d.Name()}
	}
	return nil
}

func (d dotFile) hasDestFile(src string) bool {
	if _, err := os.Stat(src); os.IsExist(err) {
		return true
	}
	return false
}

func (d dotFile) whiteListed() bool {
	return d.whitelisted
}

func (d dotFile) createRootFile() error { return nil }

func (d dotFile) Name() string {
	return d.entry.Name()
}

func (d dotFile) hasDotPrefix() bool {
	return strings.HasPrefix(d.Name(), ".")
}

type SymbolicLinkError struct {
	err   error
	fname string
}

func (e SymbolicLinkError) Error() string {
	return fmt.Sprintf("%s failed symbolic link: %v \n", e.fname, e.err)
}

func (e SymbolicLinkError) Unwrap() error {
	return e.err
}

// guarantee that a filename is unique
func unique(files ...string) []string {
	keys := make(map[string]struct{})
	var uniqueFiles []string

	for _, f := range files {
		if _, ok := keys[f]; !ok {
			keys[f] = struct{}{}
			uniqueFiles = append(uniqueFiles, f)
		}
	}
	return uniqueFiles
}
