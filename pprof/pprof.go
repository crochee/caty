// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/31

package pprof

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func NewPprof(cpuPath, memPath string) *profiler {
	if cpuPath == "" {
		cpuPath = "./"
	}
	if memPath == "" {
		memPath = "./"
	}
	return &profiler{
		cpuPath: cpuPath,
		memPath: memPath,
	}
}

type profiler struct {
	sync.Mutex
	running bool
	exit    chan bool

	// where the cpu profile is written
	cpuFile *os.File
	// where the mem profile is written
	memFile *os.File

	cpuPath string
	memPath string
}

func (p *profiler) writeHeap() {
	defer p.memFile.Close()

	t := time.NewTicker(time.Second * 30)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			runtime.GC()
			_ = pprof.WriteHeapProfile(p.memFile)
		case <-p.exit:
			return
		}
	}
}

func (p *profiler) Start() error {
	p.Lock()
	defer p.Unlock()

	if p.running {
		return nil
	}

	// create exit channel
	p.exit = make(chan bool)

	cpuFile := filepath.Clean(p.cpuPath + "cpu.pprof")
	memFile := filepath.Clean(p.memPath + "mem.pprof")

	var err error
	if p.cpuFile, err = os.Create(cpuFile); err != nil {
		return err
	}

	if p.memFile, err = os.Create(memFile); err != nil {
		return err
	}

	// start cpu profiling
	if err = pprof.StartCPUProfile(p.memFile); err != nil {
		return err
	}

	// write the heap periodically
	go p.writeHeap()

	p.running = true
	return nil
}

func (p *profiler) Stop() error {
	p.Lock()
	defer p.Unlock()

	select {
	case <-p.exit:
		return nil
	default:
		close(p.exit)
		pprof.StopCPUProfile()
		_ = p.cpuFile.Close()
		p.running = false
		p.cpuFile = nil
		p.memFile = nil
		return nil
	}
}

func (p *profiler) String() string {
	return "obs_pprof"
}
