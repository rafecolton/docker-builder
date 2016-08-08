package fileutils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// CpArgs is a list of arguments which can be passed to the CpWithArgs method
type CpArgs struct {
	Recursive       bool
	PreserveLinks   bool
	PreserveModTime bool
}

// Cp is like `cp`
func Cp(src, dest string) (err error) {
	return CpWithArgs(src, dest, CpArgs{})
}

func cpSymlink(src, dest string) (err error) {
	var linkTarget string
	linkTarget, err = os.Readlink(src)
	if err != nil {
		return
	}

	return os.Symlink(linkTarget, dest)
}

func cpFollowLinks(src, dest string) (err error) {
	return CpWithArgs(src, dest, CpArgs{})
}

func cpPreserveLinks(src, dest string) (err error) {
	return CpWithArgs(src, dest, CpArgs{PreserveLinks: true})
}

/*
CpR is like `cp -R`
*/
func CpR(source, dest string) (err error) {
	return CpWithArgs(source, dest, CpArgs{Recursive: true, PreserveLinks: true})
}

/*
CpWithArgs is a version of the Cp method that allows the passing of an
arguments struct to further modify the copying behavior
*/
func CpWithArgs(source, dest string, args CpArgs) (err error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return
	}

	if sourceInfo.IsDir() {
		// Handle the dir case
		if !args.Recursive {
			return errors.New("source is a directory")
		}

		// ensure dest dir does not already exist
		if _, err = os.Open(dest); !os.IsNotExist(err) {
			return errors.New("destination already exists")
		}

		// create dest dir
		if err = os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
			return
		}

		files, err := ioutil.ReadDir(source)
		if err != nil {
			return err
		}

		for _, file := range files {
			sourceFilePath := fmt.Sprintf("%s/%s", source, file.Name())
			destFilePath := fmt.Sprintf("%s/%s", dest, file.Name())

			if err = CpWithArgs(sourceFilePath, destFilePath, args); err != nil {
				return err
			}
		}
	} else {
		// Handle the file case
		si, err := os.Lstat(source)
		if err != nil {
			return err
		}

		if args.PreserveLinks && !si.Mode().IsRegular() {
			return cpSymlink(source, dest)
		}

		//open source
		in, err := os.Open(source)
		if err != nil {
			return err
		}
		defer in.Close()

		//create dest
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()

		//copy to dest from source
		if _, err = io.Copy(out, in); err != nil {
			return err
		}

		if err = out.Chmod(si.Mode()); err != nil {
			return err
		}

		if args.PreserveModTime {
			if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
				return err
			}
		}

		//sync dest to disk
		err = out.Sync()
	}

	return
}
