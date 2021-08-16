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

package inject

import (
	"github.com/vine-io/pkg/inject"
	log "github.com/vine-io/vine/lib/logger"
)

type logger struct{}

func (l logger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func init() {
	g = inject.Container{}
	g.Logger = logger{}
}

var g inject.Container

func Provide(vv ...interface{}) error {
	for _, v := range vv {
		if err := g.Provide(&inject.Object{
			Value: v,
		}); err != nil {
			return err
		}
	}

	return nil
}

func ProvidePanic(vv ...interface{}) {
	if err := Provide(vv...); err != nil {
		panic(err)
	}
}

func ProvideWithName(v interface{}, name string) error {
	return g.Provide(&inject.Object{Value: v, Name: name})
}

func ProvideWithNamePanic(v interface{}, name string) {
	if err := ProvideWithName(v, name); err != nil {
		panic(err)
	}
}

func PopulateTarget(target interface{}) error {
	return g.PopulateTarget(target)
}

func Populate() error {
	return g.Populate()
}

func Objects() []*inject.Object {
	return g.Objects()
}

func Resolve(dst interface{}) error {
	return g.Resolve(dst)
}

func ResolveByName(dst interface{}, name string) error {
	return g.ResolveByName(dst, name)
}
