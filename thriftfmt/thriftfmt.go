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

	"github.com/trigger3/toy/thriftfmt/file"
	"github.com/trigger3/toy/thriftfmt/formator"
	"github.com/trigger3/toy/thriftfmt/key_words"
	"github.com/trigger3/toy/thriftfmt/line_parser"
)

var (
	// main operation modes
	write  = flag.Bool("w", false, "write result to (source) file instead of stdout")
	list   = flag.Bool("l", false, "list files whose formatting differs from thriftfmt's")
	doDiff = flag.Bool("d", false, "display diffs instead of rewriting files")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: thriftfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		if *write {
			fmt.Fprintln(os.Stderr, "error: cannot use -w with standard input")
			os.Exit(-1)
			return
		}
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-2)
			return
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-3)
			return
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(-2)
				// TODO err
			}
		}
	}
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isThriftFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	// Don't complain if a file was deleted in the meantime (i.e.
	// the directory changed concurrently while running thriftfmt).
	if err != nil && !os.IsNotExist(err) {
		// TODO err
		panic(err.Error())
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func isThriftFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".thrift")
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
	i := 0
	for scanner.Scan() {
		i++
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if len(line) == 0 {
			continue
		}
		terms := lineParser.ParseOneLine(line)
		if terms == nil || len(terms) == 0 {
			continue
		}
		if err := formator1.Format(terms); err != nil {
			return nil, fmt.Errorf("%v:%v, state:%v, err:%w", file, i, line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	res := make([]byte, 0)
	buff := bytes.NewBuffer(res)
	formator1.Print(buff)

	return buff.Bytes(), nil
}
