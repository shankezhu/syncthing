// Copyright (C) 2017 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// +build !windows

package fs

import (
	"os"
	"path/filepath"
)

// DebugSymlinkForTestsOnly is not and should not be used in Syncthing code,
// hence the cumbersome name to make it obvious if this ever leaks. Its
// reason for existence is the Windows version, which allows creating
// symlinks when non-elevated.
func DebugSymlinkForTestsOnly(oldFs, newFs Filesystem, oldname, newname string) error {
	if caseFs, ok := unwrapFilesystem(newFs).(*caseFilesystem); ok {
		if err := caseFs.checkCase(newname); err != nil {
			return err
		}
		caseFs.dropCache()
	}
	if err := os.Symlink(filepath.Join(oldFs.URI(), oldname), filepath.Join(newFs.URI(), newname)); err != nil {
		return err
	}
	return nil
}

// unwrapFilesystem removes "wrapping" filesystems to expose the underlying filesystem.
func unwrapFilesystem(fs Filesystem) Filesystem {
	for {
		switch sfs := fs.(type) {
		case *logFilesystem:
			fs = sfs.Filesystem
		case *walkFilesystem:
			fs = sfs.Filesystem
		case *MtimeFS:
			fs = sfs.Filesystem
		default:
			return sfs
		}
	}
}
