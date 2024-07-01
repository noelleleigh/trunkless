package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const p = "/home/vilmibm/geocities/UPPERCASE/geocities/YAHOOIDS"

func main() {
	userDirs := make(chan string)
	var wg sync.WaitGroup

	walkFn := func(s string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		isUserDir, err := filepath.Match("/home/vilmibm/geocities/UPPERCASE/geocities/YAHOOIDS/?/?/*", s)
		if err != nil {
			return err
		}

		if isUserDir {
			userDirs <- s
		}

		//if len(d.Name()) > 1 && d.IsDir() {
		//	fmt.Printf("%s %s\n", s, d.Name())
		//}
		// TODO be able to tell when s is a full path (ie to a file for user)
		// TODO sniff what kind of file full path points to
		// TODO if text, read and append to file for that user
		return nil
	}
	go func() {
		err := filepath.WalkDir(p, walkFn)
		close(userDirs)
		if err != nil {
			panic(err)
		}
	}()
	for ud := range userDirs {
		wg.Add(1)
		go processUserDir(&wg, ud)
	}
	wg.Wait()
}

func processUserDir(wg *sync.WaitGroup, ud string) {
	defer wg.Done()
	fmt.Printf("GONNA PROCESS %s\n", ud)
}
