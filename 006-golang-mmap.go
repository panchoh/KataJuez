package main

import (
	"bytes"
	"log"
	"os"
	"syscall"
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

	mmap1, err := syscall.Mmap(int(f1.Fd()), 0, int(f1info.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Could not mmap() %v: %v", file1, err)
	}
	defer syscall.Munmap(mmap1)

	mmap2, err := syscall.Mmap(int(f2.Fd()), 0, int(f2info.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Could not mmap() %v: %v", file2, err)
	}
	defer syscall.Munmap(mmap2)

	if bytes.Equal(mmap1, mmap2) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
