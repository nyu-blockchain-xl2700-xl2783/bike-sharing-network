/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	ccmetadata "github.com/hyperledger/fabric/core/common/ccprovider/metadata"
	"github.com/pkg/errors"
)

var vmLogger = flogging.MustGetLogger("container")

// These filetypes are excluded while creating the tar package sent to Docker
// Generated .class and other temporary files can be excluded
var javaExcludeFileTypes = map[string]bool{
	".class": true,
}

// WriteFolderToTarPackage writes source files to a tarball.
// This utility is used for node js chaincode packaging, but not golang chaincode.
// Golang chaincode has more sophisticated file packaging, as implemented in golang/platform.go.
func WriteFolderToTarPackage(tw *tar.Writer, srcPath string, excludeDir string, includeFileTypeMap map[string]bool, excludeFileTypeMap map[string]bool) error {
	fileCount := 0
	rootDirectory := srcPath

	// trim trailing slash if it was passed
	if rootDirectory[len(rootDirectory)-1] == '/' {
		rootDirectory = rootDirectory[:len(rootDirectory)-1]
	}

	vmLogger.Infof("rootDirectory = %s", rootDirectory)

	//append "/" if necessary
	if excludeDir != "" && strings.LastIndex(excludeDir, "/") < len(excludeDir)-1 {
		excludeDir = excludeDir + "/"
	}

	rootDirLen := len(rootDirectory)
	walkFn := func(path string, info os.FileInfo, err error) error {

		// If path includes .git, ignore
		if strings.Contains(path, ".git") {
			return nil
		}

		if info.Mode().IsDir() {
			return nil
		}

		//exclude any files with excludeDir prefix. They should already be in the tar
		if excludeDir != "" && strings.Index(path, excludeDir) == rootDirLen+1 {
			//1 for "/"
			return nil
		}
		// Because of scoping we can reference the external rootDirectory variable
		if len(path[rootDirLen:]) == 0 {
			return nil
		}
		ext := filepath.Ext(path)

		if includeFileTypeMap != nil {
			// we only want 'fileTypes' source files at this point
			if _, ok := includeFileTypeMap[ext]; ok != true {
				return nil
			}
		}

		//exclude the given file types
		if excludeFileTypeMap != nil {
			if exclude, ok := excludeFileTypeMap[ext]; ok && exclude {
				return nil
			}
		}

		var newPath string

		// if file is metadata, keep the /META-INF directory, e.g: META-INF/statedb/couchdb/indexes/indexOwner.json
		// otherwise file is source code, put it in /src dir, e.g: src/marbles_chaincode.js
		if strings.HasPrefix(path, filepath.Join(rootDirectory, "META-INF")) {
			newPath = path[rootDirLen+1:]

			// Split the filename itself from its path
			_, filename := filepath.Split(newPath)

			// Hidden files are not supported as metadata, therefore ignore them.
			// User often doesn't know that hidden files are there, and may not be able to delete them, therefore warn user rather than error out.
			if strings.HasPrefix(filename, ".") {
				vmLogger.Warningf("Ignoring hidden file in metadata directory: %s", newPath)
				return nil
			}

			// Validate metadata file for inclusion in tar
			// Validation is based on the passed metadata directory, e.g. META-INF/statedb/couchdb/indexes
			err = ccmetadata.ValidateMetadataFile(path, filepath.Dir(newPath))
			if err != nil {
				return err
			}

		} else { // file is not metadata, include in src
			newPath = fmt.Sprintf("src%s", path[rootDirLen:])
		}

		err = WriteFileToPackage(path, newPath, tw)
		if err != nil {
			return fmt.Errorf("Error writing file to package: %s", err)
		}
		fileCount++

		return nil
	}

	if err := filepath.Walk(rootDirectory, walkFn); err != nil {
		vmLogger.Infof("Error walking rootDirectory: %s", err)
		return err
	}
	// return error if no files were found
	if fileCount == 0 {
		return errors.Errorf("no source files found in '%s'", srcPath)
	}
	return nil
}

//Package Java project to tar file from the source path
func WriteJavaProjectToPackage(tw *tar.Writer, srcPath string) error {

	vmLogger.Debugf("Packaging Java project from path %s", srcPath)

	if err := WriteFolderToTarPackage(tw, srcPath, "", nil, javaExcludeFileTypes); err != nil {

		vmLogger.Errorf("Error writing folder to tar package %s", err)
		return err
	}
	// Write the tar file out
	if err := tw.Close(); err != nil {
		return err
	}
	return nil

}

//WriteFileToPackage writes a file to the tarball
func WriteFileToPackage(localpath string, packagepath string, tw *tar.Writer) error {
	vmLogger.Debug("Writing file to tarball:", packagepath)
	fd, err := os.Open(localpath)
	if err != nil {
		return fmt.Errorf("%s: %s", localpath, err)
	}
	defer fd.Close()

	is := bufio.NewReader(fd)
	return WriteStreamToPackage(is, localpath, packagepath, tw)

}

//WriteStreamToPackage writes bytes (from a file reader) to the tarball
func WriteStreamToPackage(is io.Reader, localpath string, packagepath string, tw *tar.Writer) error {
	info, err := os.Stat(localpath)
	if err != nil {
		return fmt.Errorf("%s: %s", localpath, err)
	}
	header, err := tar.FileInfoHeader(info, localpath)
	if err != nil {
		return fmt.Errorf("Error getting FileInfoHeader: %s", err)
	}

	//Let's take the variance out of the tar, make headers identical by using zero time
	oldname := header.Name
	var zeroTime time.Time
	header.AccessTime = zeroTime
	header.ModTime = zeroTime
	header.ChangeTime = zeroTime
	header.Name = packagepath
	header.Mode = 0100644
	header.Uid = 500
	header.Gid = 500

	if err = tw.WriteHeader(header); err != nil {
		return fmt.Errorf("Error write header for (path: %s, oldname:%s,newname:%s,sz:%d) : %s", localpath, oldname, packagepath, header.Size, err)
	}
	if _, err := io.Copy(tw, is); err != nil {
		return fmt.Errorf("Error copy (path: %s, oldname:%s,newname:%s,sz:%d) : %s", localpath, oldname, packagepath, header.Size, err)
	}

	return nil
}

func WriteBytesToPackage(name string, payload []byte, tw *tar.Writer) error {
	//Make headers identical by using zero time
	var zeroTime time.Time
	tw.WriteHeader(&tar.Header{Name: name, Size: int64(len(payload)), ModTime: zeroTime, AccessTime: zeroTime, ChangeTime: zeroTime})
	tw.Write(payload)

	return nil
}
