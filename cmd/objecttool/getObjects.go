package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/masiulaniec/Dominator/lib/fsutil"
	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/objectcache"
	"github.com/masiulaniec/Dominator/lib/objectserver"
)

func getObjectsSubcommand(objSrv objectserver.ObjectServer, args []string) {
	if err := getObjects(objSrv, args[0], args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting objects: %s\n", err)
		os.Exit(2)
	}
	os.Exit(0)
}

func getObjects(objSrv objectserver.ObjectServer,
	hashesFilename, outputDirectory string) error {
	hashesFile, err := os.Open(hashesFilename)
	if err != nil {
		return err
	}
	defer hashesFile.Close()
	scanner := bufio.NewScanner(hashesFile)
	var hashes []hash.Hash
	for scanner.Scan() {
		hashval, err := objectcache.FilenameToHash(scanner.Text())
		if err != nil {
			return err
		}
		hashes = append(hashes, hashval)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	objectsReader, err := objSrv.GetObjects(hashes)
	if err != nil {
		return err
	}
	defer objectsReader.Close()
	tmpDirname := outputDirectory + "~"
	if err := os.Mkdir(tmpDirname, syscall.S_IRWXU); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDirname)
	for _, hash := range hashes {
		length, reader, err := objectsReader.NextObject()
		if err != nil {
			return err
		}
		err = readOne(tmpDirname, hash, length, reader)
		reader.Close()
		if err != nil {
			return err
		}
	}
	return os.Rename(tmpDirname, outputDirectory)
}

func readOne(dirname string, hash hash.Hash, length uint64,
	reader io.Reader) error {
	filename := fmt.Sprintf("%s/%x", dirname, hash)
	return fsutil.CopyToFile(filename, filePerms, reader, length)
}
