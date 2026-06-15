package main

import (
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"
)

var is_pprofing bool
var pp_file *os.File

func start_pprof() error {
	os.Mkdir("./pprof", os.ModeDir)

	file_name := time.Now().String()
	file_name = strings.Join(strings.Split(file_name, " "), "-") + ".prof"
	file, err := os.Create(filepath.Join("./pprof", file_name))
	if err != nil {
		return err
	}

	is_pprofing = true
	pp_file = file
	return pprof.StartCPUProfile(file)
}

func stop_pprof() {
	pprof.StopCPUProfile()
	pp_file.Close()
	is_pprofing = false
}
