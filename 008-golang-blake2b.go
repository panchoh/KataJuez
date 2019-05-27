package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/blake2b"
)

func main() {
	file1 := os.Args[1]

	f1, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open %v: %v", file1, err)
	}
	defer func() {
		if err := f1.Close(); err != nil {
			log.Panic(err)
		}
	}()

	f1info, err := f1.Stat()
	if err != nil {
		log.Fatalf("could not Stat() %v: %v", file1, err)
	}

	file2 := os.Args[2]

	f2, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatalf("could not open %v: %v", file2, err)
	}
	defer func() {
		if err := f2.Close(); err != nil {
			log.Panic(err)
		}
	}()

	f2info, err := f2.Stat()
	if err != nil {
		log.Fatalf("could not Stat() %v: %v", file2, err)
	}

	if f1info.Size() != f2info.Size() {
		os.Exit(1)
	}

	buf1 := make([]byte, 8192)
	buf2 := make([]byte, 8192)

	hasher1, _ := blake2b.New512(nil)
	hasher2, _ := blake2b.New512(nil)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			n1, err := f1.Read(buf1)
			if err != nil && err != io.EOF {
				log.Fatalf("error reading %v: %v", file1, err)
			}
			if n1 == 0 {
				return
			}
			hasher1.Write(buf1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			n2, err := f2.Read(buf2)
			if err != nil && err != io.EOF {
				log.Fatalf("error reading %v: %v", file2, err)
			}
			if n2 == 0 {
				return
			}
			hasher2.Write(buf2)
		}
	}()

	wg.Wait()

	if bytes.Compare(hasher1.Sum(nil), hasher2.Sum(nil)) != 0 {
		os.Exit(1)
	}
}
