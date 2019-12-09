// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bingoohuang/gostarter/model"

	"github.com/gin-gonic/gin"
)

// GinRecovery wraps the structure of gin recovery info.
type GinRecovery struct {
	dunno     []byte
	centerDot []byte
	dot       []byte
	slash     []byte
	reset     string
}

// MakeGinRecovery makes the GinRecovery
func MakeGinRecovery() *GinRecovery {
	return &GinRecovery{
		dunno:     []byte("???"),
		centerDot: []byte("·"),
		dot:       []byte("."),
		slash:     []byte("/"),
		reset:     string([]byte{27, 91, 48, 109}),
	}
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func (g *GinRecovery) Recovery() gin.HandlerFunc {
	return g.RecoveryWithWriter(gin.DefaultErrorWriter)
}

// RecoveryWithWriter returns a middleware for a given writer
// that recovers from any panics and writes a 500 if there was one.
func (g *GinRecovery) RecoveryWithWriter(out io.Writer) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				g.dealRecover(err, logger, c)
			}
		}()

		c.Next()
	}
}

func (g *GinRecovery) dealRecover(err interface{}, logger *log.Logger, c *gin.Context) {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool

	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			e := strings.ToLower(se.Error())
			if strings.Contains(e, "broken pipe") ||
				strings.Contains(e, "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	if logger != nil {
		stack := g.stack(3)
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		headers := strings.Split(string(httpRequest), "\r\n")

		for idx, header := range headers {
			current := strings.Split(header, ":")
			if current[0] == "Authorization" {
				headers[idx] = current[0] + ": *"
			}
		}

		switch {
		case brokenPipe:
			logger.Printf("%s\n%s%s", err, string(httpRequest), g.reset)
		case gin.IsDebugging():
			logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
				g.timeFormat(time.Now()), strings.Join(headers, "\r\n"), err, stack, g.reset)
		default:
			logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
				g.timeFormat(time.Now()), err, stack, g.reset)
		}
	}

	// If the connection is dead, we can't write a status to it.
	if brokenPipe {
		c.Error(err.(error)) // nolint: errcheck
		c.Abort()
	} else {
		rsp := model.Rsp{
			Status:  500,
			Message: "系统异常",
		}
		c.JSON(http.StatusOK, rsp)
		c.Abort()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func (g *GinRecovery) stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data

	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte

	var lastFile string

	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		fmt.Fprintf(buf, "\t%s: %s\n", g.function(pc), g.source(lines, line))
	}

	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func (g *GinRecovery) source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed

	if n < 0 || n >= len(lines) {
		return g.dunno
	}

	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func (g *GinRecovery) function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return g.dunno
	}

	name := []byte(fn.Name())

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, g.slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	if period := bytes.Index(name, g.dot); period >= 0 {
		name = name[period+1:]
	}

	return bytes.Replace(name, g.centerDot, g.dot, -1)
}

func (g *GinRecovery) timeFormat(t time.Time) string {
	return t.Format("2006/01/02 - 15:04:05")
}
