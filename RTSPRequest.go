/**
* The MIT License (MIT)
*
* Copyright (c) 2016 doublemo<435420057@qq.com>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
*/

package koala

import (
    "bufio"
    "net/http"
    "net/url"
    "io"
    "io/ioutil"
    "fmt"
    "strconv"
    "strings"
)

type IRTSPRequest interface {
    GetMethod()  string

    GetURL()    *url.URL

    GetVersion() string

    GetHeader()  http.Header

    GetContentLength() int

    GetBody()    io.ReadCloser
}

type RTSPRequest struct {
    Method      string
    URL         *url.URL
    Proto        string
    ProtoMajor    int
    ProtoMinor    int
    Header        http.Header
    ContentLength int
    Body          io.ReadCloser
}

func (request RTSPRequest) String() string {
    s := fmt.Sprintf("%s %s %s/%d.%d\r\n", request.Method, request.URL, request.Proto, request.ProtoMajor, request.ProtoMinor)
    for k, v := range request.Header {
        for _, v := range v {
            s += fmt.Sprintf("%s: %s\r\n", k, v)
        }
    }
    s += "\r\n"
    if request.Body != nil {
        str, _ := ioutil.ReadAll(request.Body)
        s += string(str)
    }
    return s
}

func (request RTSPRequest) parseVersion( s string ) (proto string, major int, minor int, err error) {
    s = strings.TrimRight(s, "\r\n")
    parts := strings.SplitN(s, "/", 2)
    proto = parts[0]
    parts = strings.SplitN(parts[1], ".", 2)
    if major, err = strconv.Atoi(parts[0]); err != nil {
        return
    }
    if minor, err = strconv.Atoi(parts[1]); err != nil {
        return
    }
    return
}

func (request *RTSPRequest) FromBytes( r io.Reader ) error {
    request.Header = make(map[string][]string)
    b      := bufio.NewReader(r)
    s, err := b.ReadString('\n')
    if err != nil {
        return err
    }

    parts           := strings.SplitN(s, " ", 3)
    request.Method   = parts[0]
    request.URL, err = url.Parse(parts[1])
    if err != nil {
        return err
    }

    request.Proto, request.ProtoMajor, request.ProtoMinor, err = request.parseVersion(parts[2])
    if err != nil {
        return err
    }

    // Parse Header
    for {
        if s, err = b.ReadString('\n'); err != nil {
            return err
        } else if s = strings.TrimRight(s, "\r\n"); s == "" {
            break
        }

        head := strings.SplitN(s, ":", 2)
        request.Header.Add(strings.TrimSpace(head[0]), strings.TrimSpace(head[1]))
    }

    request.ContentLength, _ = strconv.Atoi(request.Header.Get("Content-Length"))
    request.Body = RequestCloser{b, r}
    return nil
}



type RequestCloser struct {
    *bufio.Reader
    r io.Reader
}

func (requestCloser RequestCloser) Close() error{
    defer func(){
        requestCloser.Reader = nil
        requestCloser.r      = nil
    }()

    if requestCloser.Reader == nil {
        return nil
    }

    if r, ok := requestCloser.r.(io.ReadCloser); ok {
        return r.Close()
    }

    return nil
}
