package main

import (
	"image/jpeg"
	"os"

	"github.com/agatan/bktree"
	"github.com/corona10/goimagehash"
)

func hashImage(image string) uint64 {
	file, _ := os.Open(image)
	defer file.Close()

	img, _ := jpeg.Decode(file)

	hash, _ := goimagehash.DifferenceHash(img)
	return hash.GetHash()
}

type image struct {
	path string
	hash uint64
}

func (x image) Distance(e bktree.Entry) int {
	count := 0
	var k uint64 = 1
	a := x.hash
	b := e.(image).hash
	for i := 0; i < 64; i++ {
		if a&k != b&k {
			count++
		}
		k <<= 1
	}
	return count
}
