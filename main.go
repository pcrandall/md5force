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

var CHARS = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var passwords []string

func main() {
	if loadHashes() {
		return
	}
	fmt.Println("File with passwords loaded. We're gonna crack", len(passwords), "passwords!")
	start := time.Now()

	cont := 2
	for len(passwords) > 0 {
		fmt.Println("Searching for passwords at length: ", cont)
		var wg sync.WaitGroup
		wg.Add(1)
		go compute(0, cont, "", &wg)
		wg.Wait()
		cont++
	}

	elapsed := time.Since(start)
	fmt.Println("Password's file cracked in:", elapsed)

}

func compute(prefix int, n int, a string, wgFather *sync.WaitGroup) {
	defer wgFather.Done()

	if prefix == n-1 {
		for _, d := range CHARS {
			password := fmt.Sprintf("%s%c", a, d)

			if searchPassword(password) {
				if len(passwords) == 0 {
					return
				}
			}
		}
	} else {
		for i := range CHARS {
			wgFather.Add(1)
			if prefix == 0 {
				go compute(prefix+1, n, fmt.Sprintf("%s%c", a, CHARS[i]), wgFather)
			} else {
				compute(prefix+1, n, fmt.Sprintf("%s%c", a, CHARS[i]), wgFather)
			}
		}
	}
}

func searchPassword(pass string) bool {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(pass)))
	for i, value := range passwords {
		if strings.Compare(hash, value) == 0 {
			// Password found!
			fmt.Println("Find Password:", pass, " with hash:", hash)
			passwords = append(passwords[:i], passwords[i+1:]...)
			return true
		}
	}
	return false
}

func loadHashes() bool {
	stream, err := ioutil.ReadFile("file.txt")
	if err != nil {
		log.Fatal(err)
		return true
	}
	readstring := string(stream)

	passwords = strings.Split(readstring, "\n")
	return false
}
