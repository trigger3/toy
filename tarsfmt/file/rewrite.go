package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func Write(filename string, src, new []byte, perm os.FileMode) error {
	// make a temporary backup before overwriting original
	bakname, err := backupFile(filename+".", src, perm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, new, perm)
	if err != nil {
		os.Rename(bakname, filename)
		return err
	}
	err = os.Remove(bakname)
	if err != nil {
		return err
	}

	return nil
}

const chmodSupported = runtime.GOOS != "windows"

func backupFile(filename string, data []byte, perm os.FileMode) (string, error) {
	// create backup file
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return "", err
	}
	bakname := f.Name()
	if chmodSupported {
		err = f.Chmod(perm)
		if err != nil {
			f.Close()
			os.Remove(bakname)
			return bakname, err
		}
	}

	// write data to backup file
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return bakname, err
}
