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

package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vine-io/vine/core/server/grpc"
)

//go:embed ca.pem client-key.pem client.pem server-key.pem server.pem
var f embed.FS

func GetSSL(root string) (*grpc.Grpc2Http, error) {
	dir := filepath.Join(root, "ssl")
	_ = os.MkdirAll(dir, 0o777)

	fca := filepath.Join(dir, "ca.pem")
	fpem := filepath.Join(dir, "server.pem")
	fkey := filepath.Join(dir, "server-key.pem")

	if isNotExists(fca) {
		ca, err := f.ReadFile("ca.pem")
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(fca, ca, 0o777); err != nil {
			return nil, err
		}
	}
	if isNotExists(fpem) {
		pem, err := f.ReadFile("server.pem")
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(fpem, pem, 0o777); err != nil {
			return nil, err
		}
	}
	if isNotExists(fkey) {
		key, err := f.ReadFile("server-key.pem")
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(fkey, key, 0o777); err != nil {
			return nil, err
		}
	}

	gh := &grpc.Grpc2Http{
		CertFile: fpem,
		KeyFile:  fkey,
		CaFile:   fca,
	}

	return gh, nil
}

func isNotExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

func GetTLS() (*tls.Config, error) {
	ca, err := f.ReadFile("ca.pem")
	if err != nil {
		return nil, err
	}
	pem, err := f.ReadFile("client.pem")
	if err != nil {
		return nil, err
	}
	key, err := f.ReadFile("client-key.pem")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(pem, key)
	if err != nil {
		return nil, err
	}

	cpool := x509.NewCertPool()
	if ok := cpool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("AppendCertFromPEM failed")
	}

	c := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      cpool,
	}

	return c, nil
}
