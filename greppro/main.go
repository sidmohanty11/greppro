package main

import (
	"fmt"
	"greppro/worker"
	"greppro/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

// > It takes a worklist and a path, and adds all the files in the path to the worklist
func discoverDirs(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("ReadDir error,", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDirs(wl, nextPath)
		} else {
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

// A struct that is used to parse the command line arguments.
var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

	var workersWg sync.WaitGroup

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := 10

	workersWg.Add(1)

// A goroutine that is responsible for discovering all the directories and adding them to the worklist.
	go func() {
		defer workersWg.Done()
		discoverDirs(&wl, args.SearchDir)
		wl.Finalize(numWorkers)
	}()

	// Creating 10 workers and adding them to the workersWg.
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		// A goroutine that is responsible for taking the next job from the worklist and passing it to the
		// worker.
		go func() {
			defer workersWg.Done()
			for {
				workEntry := wl.Next()
				if workEntry.Path != "" {
					workerRes := worker.FindInFile(workEntry.Path, args.SearchTerm)
					if workerRes != nil {
						for _, r := range workerRes.Inner {
							results <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})

// Waiting for the workers to finish and then closing the channel.
	go func() {
		workersWg.Wait()
		close(blockWorkersWg)
	}()

	var displayWg sync.WaitGroup

	displayWg.Add(1)

	// A goroutine that is responsible for printing the results.
	go func() {
		for {
			select {
			case r := <- results:
				fmt.Println(r.Path, r.LineNum, r.Line)
			case <-blockWorkersWg:
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			}
		}
	}()

	displayWg.Wait()
}