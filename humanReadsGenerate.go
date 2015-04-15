package main

import (
	"os"
	"bufio"
	"fmt"
	"io/ioutil"
	"bytes"
	"strconv"
	"math/rand"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		panic("Must provide sequence folder, kmer length and reads number!")
	}

	files, _ := ioutil.ReadDir(os.Args[1])
	K, _ := strconv.Atoi(os.Args[2])
	N,_ := strconv.Atoi(os.Args[3])
	
	contigs := make([]byte, 0)
	for _, file := range files {
		contig := contigExtract(os.Args[1]+file.Name())
		contigs = append(contigs, contig...)
	}

	RC := map[byte]byte {'A': 'T', 'T': 'A', 'G': 'C', 'C': 'G'}

	rand.Seed(time.Now().UTC().UnixNano())
	n := 0
	for n<N {
		strand := rand.Intn(2)
		index := rand.Intn(len(contigs)-K)
		if !bytes.Contains(contigs[index:index+K], []byte("N")) {
			if strand == 0 {
				fmt.Println(string(contigs[index:index+K]))
			} else {
				reverse_read := make([]byte, K)
				for k, c := range contigs[index:index+K] {
					reverse_read[k] = RC[c]
				}
				fmt.Println(string(reverse_read))
			}
			n++
		}
	}
}

func contigExtract(filename string) ([]byte) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer f.Close()
	br := bufio.NewReader(f)
	byte_buffer := bytes.Buffer{}
	_, isPrefix, err := br.ReadLine()
	if err != nil || isPrefix {
		fmt.Printf("%v\n", nil)
		os.Exit(1)
	}

	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil || isPrefix {
			byte_buffer.Write([]byte("N"))
			break
		} else {
			if bytes.Contains(line, []byte(">")) {
				byte_buffer.Write([]byte("N"))
			} else {
				byte_buffer.Write(line)
			}
		}
	}
	return byte_buffer.Bytes()
}