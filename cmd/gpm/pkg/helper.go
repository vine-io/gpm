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
	"strings"

	"github.com/vine-io/cli"
	vclient "github.com/vine-io/vine/core/client"
)

func getCallOptions(c *cli.Context) []vclient.CallOption {
	host := c.String("host")
	if index := strings.Index(host, ":"); index == -1 {
		host += ":7700"
	}
	return []vclient.CallOption{
		vclient.WithAddress(),
		vclient.WithDialTimeout(c.Duration("dial-timeout")),
		vclient.WithRequestTimeout(c.Duration("request-timeout")),
		vclient.WithStreamTimeout(c.Duration("request-timeout")),
	}
}
