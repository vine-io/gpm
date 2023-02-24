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

// Code generated by vine. DO NOT EDIT.
package internal

const Namespace = "io.gpm"

const GpmName = "io.gpm.service.gpm"
const GpmId = "12df4108-6176-4c93-ac58-bfa003d4ddb1"

var (
	GitTag     = "v1.10.2"
	GitCommit  = "d9eeb18"
	BuildDate  = "1677229852"
	GetVersion = func() string {
		v := GitTag
		if GitCommit != "" {
			v += "-" + GitCommit
		}
		if BuildDate != "" {
			v += "-" + BuildDate
		}
		if v == "" {
			return "latest"
		}
		return v
	}
)
