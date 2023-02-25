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

package ctl

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func tarFn(c *cobra.Command, args []string) error {

	outE := os.Stdout
	name, _ := c.Flags().GetString("name")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	targets, _ := c.Flags().GetStringSlice("target")
	if len(targets) == 0 {
		return fmt.Errorf("missing target")
	}

	fmt.Fprintf(outE, "starting tar %s\n", name)
	dest, err := os.Create(name)
	if err != nil {
		return err
	}
	defer dest.Close()

	gw := gzip.NewWriter(dest)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, t := range targets {
		err = compress(t, "", tw, outE)
		if err != nil {
			return fmt.Errorf("compress %s: %v", t, err)
		}
	}

	fmt.Fprintf(outE, "tar %s successfully!\n", name)

	return nil
}

func compress(path string, prefix string, tw *tar.Writer, out io.Writer) error {
	fmt.Fprintf(out, "compress %s\n", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = filepath.Join(prefix, info.Name())
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			err = compress(filepath.Join(file.Name(), fi.Name()), prefix, tw, out)
			if err != nil {
				return err
			}
		}
		return nil
	}
	header, err := tar.FileInfoHeader(info, "")
	header.Name = filepath.Join(prefix, header.Name)
	if err != nil {
		return err
	}
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}

func TarCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tar",
		Short:   "create a compress package for Install subcommand",
		GroupID: "gpm",
		RunE:    tarFn,
	}

	cmd.PersistentFlags().StringP("name", "N", "", "the specify the name for package")
	cmd.PersistentFlags().StringSliceP("target", "T", []string{}, "the specify the target list for package")

	return cmd
}
