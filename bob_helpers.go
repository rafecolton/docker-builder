package bob

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/*
CopyFile copies one file from source to dest.  Copied from
https://github.com/opesun/copyrecur.
*/
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return
}

/*
CopyDir recursively copies one dir from source to dest.  Copied from
https://github.com/opesun/copyrecur.
*/
func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("Source is not a directory")
	}
	// ensure dest dir does not already exist
	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return errors.New("Destination already exists")
	}
	// create dest dir
	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return
}
