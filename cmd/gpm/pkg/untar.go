// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vine-io/cli"
)

func unTarFn(c *cli.Context) error {

	outE := os.Stdout
	pack := c.String("package")
	if len(pack) == 0 {
		return fmt.Errorf("missing package")
	}
	target := c.String("target")
	if len(target) == 0 {
		return fmt.Errorf("missing target")
	}

	fmt.Fprintf(outE, "starting decompress tar %s\n", pack)
	dest, err := os.Open(pack)
	if err != nil {
		return err
	}
	defer dest.Close()

	gr, err := gzip.NewReader(dest)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	defer gr.Close()

	if err = decompress(target, tr, outE); err != nil {
		return err
	}

	fmt.Fprintf(outE, "decompress %s successfully!\n", pack)

	return nil
}

func decompress(dir string, tr *tar.Reader, out io.Writer) error {
	for {
		hdr, e := tr.Next()
		if e != nil {
			if e == io.EOF {
				break
			} else {
				return e
			}
		}
		if hdr.FileInfo().IsDir() {
			_ = os.MkdirAll(filepath.Join(dir, hdr.Name), os.ModePerm)
		} else {
			fname := filepath.Join(dir, hdr.Name)
			fmt.Fprintf(out, "%s\n", fname)
			f, e1 := createFile(fname)
			if e1 != nil {
				return e1
			}
			_, e1 = io.Copy(f, tr)
			if e1 != nil && e1 != io.EOF {
				f.Close()
				return e1
			}
			f.Close()
		}
	}

	return nil
}

func UnTarCmd() *cli.Command {
	return &cli.Command{
		Name:   "untar",
		Usage:  "decompress a package",
		Action: unTarFn,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"P"},
				Usage:   "the specify the package for compressing",
			},
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"C"},
				Usage:   "the specify the target for saving file",
				Value:   ".",
			},
		},
	}
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0o755)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
}
