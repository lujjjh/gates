package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/gates/gates"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var timelimit = flag.Int("timelimit", 0, "max time to run (in seconds)")

func readSource(filename string) ([]byte, error) {
	if filename == "" || filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

func run() (gates.Value, error) {
	filename := flag.Arg(0)
	src, err := readSource(filename)
	if err != nil {
		return nil, err
	}

	if filename == "" || filename == "-" {
		filename = "<stdin>"
	}

	vm := gates.New()

	ctx := context.Background()
	if *timelimit > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*timelimit)*time.Second)
		defer cancel()
	}

	prg, err := gates.Compile(string(src))
	if err != nil {
		return nil, err
	}
	return vm.RunProgram(ctx, prg)
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			debug.Stack()
			panic(x)
		}
	}()
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	v, err := run()
	if err != nil {
		log.Println(err)
		os.Exit(64)
	}
	fmt.Println(v.ToString())
}
