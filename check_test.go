package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestDircheck(t *testing.T) {
	goodRoot := "/home/user"
	goodDest := "/home/user/Downloads"
	goodDotDotDest := "/home/user/../../home/user/Downloads"
	badDest := "/usr/etc/shadow"
	badDotDotDest := "/home/user/../../etc/../etc/shadow"

	if !Dircheck(goodRoot, goodDest) {
		t.Error("Dircheck failed on good input")
	}

	if !Dircheck(goodRoot, goodDotDotDest) {
		t.Error("Dircheck failed on good input")
	}

	if Dircheck(goodRoot, badDest) {
		t.Error("Dircheck passed on bad input")
	}

	if Dircheck(goodRoot, badDotDotDest) {
		t.Error("Dircheck passed on bad input")
	}
}

func TestDotdotcheck(t *testing.T) {
	goodDest := "/home/user/Downloads"
	goodDotDotDest := "/home/user/../../home/user/Downloads"
	badDest := "/usr/etc/shadow"
	badDotDotDest := "/home/user/../../etc/../etc/shadow"

	if !Dotdotcheck(goodDest) {
		t.Error("Dotdotcheck failed on good input")
	}

	if Dotdotcheck(goodDotDotDest) {
		t.Error("Dotdotcheck passed on bad input")
	}

	if !Dotdotcheck(badDest) {
		t.Error("Dotdotcheck failed on good input")
	}

	if Dotdotcheck(badDotDotDest) {
		t.Error("Dotdotcheck failed on bad input")
	}
}

func TestTypecheck(t *testing.T) {
	fs := afero.NewOsFs()
	fakeFilesystem := &afero.Afero{Fs: fs}
	normalFile, err := fakeFilesystem.TempFile("", "normalFile")
	if err != nil {
		t.Error("Failed to create normal file")
	}
	normalDir, err := fakeFilesystem.TempDir("", "normalDir")
	if err != nil {
		t.Error("Failed to create normal directory")
	}
	f, ok := fs.(*afero.OsFs)
	if !ok {
		t.Error("Failed to cast to OsFs")
	}
	err = f.SymlinkIfPossible(normalFile.Name(), "symlinkFile")
	if err != nil {
		t.Error("Failed to create symlink")
	}

	s, _, err := f.LstatIfPossible(normalFile.Name())
	if !Typecheck(s) {
		t.Error("Typecheck failed on good input")
	}
	s, _, err = f.LstatIfPossible(normalDir)
	if !Typecheck(s) {
		t.Error("Typecheck failed on good input")
	}
	s, _, err = f.LstatIfPossible("symlinkFile")
	if Typecheck(s) {
		t.Error("Typecheck passed on bad input")
	}

	// Now clean up the temp files
	fakeFilesystem.Remove(normalFile.Name())
	fakeFilesystem.Remove(normalDir)
	fakeFilesystem.Remove("symlinkFile")
}

func TestSymlinkSafetyCheck(t *testing.T) {
	fs := afero.NewOsFs()
	rootDir := "/home/user/Downloads"
	f := fs.(*afero.OsFs)
	err := f.SymlinkIfPossible("/home/user/Downloads/something.txt", "goodSymlinkFile")
	if err != nil {
		t.Error("Failed to create symlink")
	}
	err = f.SymlinkIfPossible("../../../../../../../../../home/user/Downloads/foo.txt", "goodDotDotSymlinkFile")
	if err != nil {
		t.Error("Failed to create symlink")
	}
	err = f.SymlinkIfPossible("/etc/shadow", "badSymlinkFile")
	if err != nil {
		t.Error("Failed to create symlink")
	}
	err = f.SymlinkIfPossible("/home/user/../../etc/shadow", "badDotDotSymlinkFile")
	if err != nil {
		t.Error("Failed to create symlink")
	}

	if !Symlinksafetycheck(rootDir, "goodSymlinkFile") {
		t.Error("Symlinksafetycheck failed on good input")
	}
	if !Symlinksafetycheck(rootDir, "goodDotDotSymlinkFile") {
		t.Error("Symlinksafetycheck failed on good input")
	}
	if Symlinksafetycheck(rootDir, "badSymlinkFile") {
		t.Error("Symlinksafetycheck passed on bad input")
	}
	if Symlinksafetycheck(rootDir, "badDotDotSymlinkFile") {
		t.Error("Symlinksafetycheck passed on bad input")
	}

	// Now clean up the temp files
	fs.Remove("goodSymlinkFile")
	fs.Remove("goodDotDotSymlinkFile")
	fs.Remove("badSymlinkFile")
	fs.Remove("badDotDotSymlinkFile")
}
