package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

var CHARS = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var HASHES []string
var wg sync.WaitGroup

func loadHashes() bool {
	stream, err := ioutil.ReadFile("file.txt")
	if err != nil {
		log.Fatal(err)
		return true
	}
	readstring := string(stream)
	HASHES = strings.Split(readstring, "\n")
	return false
}

func buildPassword(prefix int, _passLen int, _passBuild string, wgFather *sync.WaitGroup) {
	defer wgFather.Done()
	if prefix == _passLen-1 {
		for _, char := range CHARS {
			password := fmt.Sprintf("%s%c", _passBuild, char)
			if searchPassword(password) {
				if len(HASHES) == 0 {
					return
				}
			}
		}
	} else {
		for i := range CHARS {
			wgFather.Add(1)
			if prefix == 0 {
				go buildPassword(prefix+1, _passLen, fmt.Sprintf("%s%c", _passBuild, CHARS[i]), wgFather)
			} else {
				buildPassword(prefix+1, _passLen, fmt.Sprintf("%s%c", _passBuild, CHARS[i]), wgFather)
			}
		}
	}
}

func searchPassword(pass string) bool {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(pass)))
	for i, value := range HASHES {
		if strings.Compare(hash, value) == 0 {
			// Password found!
			fmt.Println("Find Password:", pass, " with hash:", hash)
			HASHES = append(HASHES[:i], HASHES[i+1:]...)
			return true
		}
	}
	return false
}

func main() {
	if loadHashes() {
		return
	}
	fmt.Println("File with hashes loaded. We're gonna crack", len(HASHES)-1, "passwords!")
	start := time.Now()
	testLen := 2
	for len(HASHES) > 1 {
		fmt.Println("Searching for passwords at length: ", testLen)
		wg.Add(1)
		go buildPassword(0, testLen, "", &wg)
		wg.Wait()
		testLen++
	}
	elapsed := time.Since(start)
	fmt.Println("Password's file cracked in:", elapsed)
}
