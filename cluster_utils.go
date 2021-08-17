package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/agatan/bktree"
	"github.com/cheggaaa/pb"
)

var clusters map[int][]string
var seen map[string]bool

func createClusters(tree bktree.BKTree, images []string, threshold int) map[int][]string {
	clusters = make(map[int][]string, len(images))
	seen = make(map[string]bool, len(images))
	var image_path string
	bar := pb.StartNew(len(images))
	for i, v := range images {
		bar.Increment()
		if !seen[v] {
			results := tree.Search(image{v, hashImage(v)}, int(threshold))
			seen[v] = true
			clusters[i] = append(clusters[i], v)
			for _, result := range results {
				image_path = string(result.Entry.(image).path)
				if !seen[image_path] {
					seen[image_path] = true
					clusters[i] = append(clusters[i], image_path)
				}
			}
		}
	}
	bar.Finish()
	fmt.Printf("Found %d clusters in %d images\n", len(clusters), len(images))
	return clusters
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func createDirectories(clusters map[int][]string, clusters_path string) {
	bar := pb.StartNew(len(clusters))
	count := 0
	for _, related_images := range clusters {
		bar.Increment()
		if len(related_images) > 1 {
			seed_path := filepath.Join(clusters_path, fmt.Sprint(count))
			err := os.Mkdir(seed_path, 0755)
			if err != nil {
				log.Fatal(err)
			}
			for _, image := range related_images {
				copy(image, filepath.Join(seed_path, filepath.Base(image)))
			}
			count += 1
		} else {
			path := filepath.Join(clusters_path, "non_repeated")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err := os.Mkdir(path, 0755)
				if err != nil {
					log.Fatal(err)
				}
			}
			copy(related_images[0], filepath.Join(path, filepath.Base(related_images[0])))
		}
	}
	bar.Finish()
}
