package bob

import (
	//"errors"
	//"fmt"
	"io"
	//"io/ioutil"
	"os"
	"os/exec"
)

/*
CopyFile copies one file from source to dest.  Copied from
https://gist.github.com/elazarl/5507969 and modified.
*/
func CopyFile(src, dest string) (err error) {
	//open source
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	//create dest
	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	//copy to dest from source
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	//duplicate source permissions on dest
	si, err := os.Stat(src)
	if err != nil {
		return
	}

	if err = out.Chmod(si.Mode()); err != nil {
		return
	}

	//sync dest to disk
	err = out.Sync()

	return
}

/*
CopyDir recursively copies one dir from source to dest.  Copied from
https://github.com/opesun/copyrecur.
*/
func CopyDir(source string, dest string) (err error) {
	return exec.Command("cp", "-af", source, dest).Run()

	/*
		THE CODE BELOW IS BROKEN - FIX IT!
	*/

	//// get properties of source dir
	//sourceInfo, err := os.Stat(source)
	//if err != nil {
	//return
	//}

	//if !sourceInfo.IsDir() {
	//return errors.New("source is not a directory")
	//}

	//// ensure dest dir does not already exist
	//if _, err = os.Open(dest); !os.IsNotExist(err) {
	//return errors.New("destination already exists")
	//}
	//// create dest dir
	//if err = os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
	//return
	//}

	//files, err := ioutil.ReadDir(source)

	//for _, file := range files {
	//sourceFilePath := fmt.Sprintf("%s/%s", source, file.Name())
	//destFilePath := fmt.Sprintf("%s/%s", dest, file.Name())

	//if file.IsDir() {
	//if err = CopyDir(sourceFilePath, destFilePath); err != nil {
	//return
	//}
	//} else {
	//if err = CopyFile(sourceFilePath, destFilePath); err != nil {
	//return
	//}
	//}

	//}
	//return
}
