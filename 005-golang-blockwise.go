package main

import (
	"bytes"
	"io"
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

	var stat syscall.Statfs_t
	err = syscall.Fstatfs(int(f1.Fd()), &stat)
	if err != nil {
		log.Fatalf("Could not determine the filesystem block size: %v", err)
	}
	blocksize := stat.Bsize

	buf1 := make([]byte, blocksize)
	buf2 := make([]byte, blocksize)

	for {
		n1, err := f1.Read(buf1)
		if err != nil && err != io.EOF {
			log.Fatalf("error reading %v: %v", file1, err)
		}
		if n1 == 0 {
			os.Exit(0)
		}

		n2, err := f2.Read(buf2)
		if err != nil && err != io.EOF {
			log.Fatalf("error reading %v: %v", file2, err)
		}
		if n2 == 0 {
			// Shouldn't happen!
			os.Exit(0)
		}

		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			os.Exit(1)
		}
	}
}
