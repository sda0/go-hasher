package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	//more efficient than crypto/sha256
	"github.com/minio/sha256-simd"
)

func main() {

	algos := [...]string{"md5", "sha1", "sha256", "sha512"}
	flags := make([]*bool, len(algos))
	for i, a := range algos {
		flags[i] = flag.Bool(a, false, fmt.Sprintf("-%s calculate %s hash of file", a, a))
	}
	flag.Parse()

	// Blocks of var/const work just as well inside a function call
	var (
		writers []io.Writer
		hashes  []hash.Hash
		names   []string
	)

	push := func(name string, h hash.Hash) {
		writers = append(writers, h)
		hashes = append(hashes, h)
		names = append(names, name)
	}

	hasher := [...]hash.Hash{md5.New(), sha1.New(), sha256.New(), sha512.New()}

	for i, flag2 := range flags {
		if *flag2 {
			push(algos[i], hasher[i])
		}
	}

	if len(names) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	in, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	io.Copy(io.MultiWriter(writers...), in)

	for i, name := range names {
		fmt.Printf("%9s: %x\n", name, hashes[i].Sum(nil))
	}
}