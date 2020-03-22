package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"toy/tarsfmt/file"
	"toy/tarsfmt/formator"
	"toy/tarsfmt/key_words"
	"toy/tarsfmt/line_parser"
)

var (
	// main operation modes
	write  = flag.Bool("w", false, "write result to (source) file instead of stdout")
	list   = flag.Bool("l", false, "list files whose formatting differs from tarsfmt's")
	doDiff = flag.Bool("d", false, "display diffs instead of rewriting files")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: tarsfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		if *write {
			fmt.Fprintln(os.Stderr, "error: cannot use -w with standard input")
			return
		}
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			fmt.Fprintln(os.Stderr, "error: %v", err)
			return
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			fmt.Fprintln(os.Stderr, "error: err:"+err.Error())
			return
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				// TODO err
				fmt.Fprintln(os.Stderr, "error: err:"+err.Error())
			}
		}
	}
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isTarsFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	// Don't complain if a file was deleted in the meantime (i.e.
	// the directory changed concurrently while running tarsfmt).
	if err != nil && !os.IsNotExist(err) {
		// TODO err
		panic(err.Error())
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func isTarsFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".tars")
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	var perm os.FileMode = 0644
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		in = f
		perm = fi.Mode().Perm()
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	res, err := format(filename, src)
	if err != nil {
		return err
	}
	if !bytes.Equal(src, res) {
		// formatting has changed
		if *list {
			fmt.Fprintln(out, filename)
		}
		if *write {
			if err := file.Write(filename, src, res, perm); err != nil {
				return err
			}
		}
		if *doDiff {
			data, err := file.Diff(src, res, filename)
			if err != nil {
				return fmt.Errorf("computing diff: %s", err)
			}
			fmt.Printf("diff -u %s %s\n", filepath.ToSlash(filename+".orig"), filepath.ToSlash(filename))
			out.Write(data)
		}
	}

	if !*list && !*write && !*doDiff {
		_, err = out.Write(res)
	}

	return err
}

func format(file string, src []byte) ([]byte, error) {
	keyWordsMgr := key_words.NewKeyWordsMgr()
	lineParser := line_parser.NewLineParse(keyWordsMgr)

	formator1 := formator.NewFormator()
	rBuff := bytes.NewBuffer(src)
	scanner := bufio.NewScanner(rBuff)
	i := 1
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if len(line) == 0 {
			i++
			continue
		}
		terms := lineParser.ParseOneLine(line)
		if err := formator1.Format(terms); err != nil {
			return nil, fmt.Errorf("%v:%v, state:%v, err:%w", file, i, line, err)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	res := make([]byte, 0)
	buff := bytes.NewBuffer(res)
	formator1.Print(buff)

	return buff.Bytes(), nil
}
