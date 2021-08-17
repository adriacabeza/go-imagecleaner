package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/agatan/bktree"
	"github.com/cheggaaa/pb"
)

var images []string

func isValidImagePath(path string) bool {
	return strings.HasSuffix(path, "jpg") || strings.HasSuffix(path, "jpeg") || strings.HasSuffix(path, "png")
}

func traverseDirectory(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	if !info.IsDir() && isValidImagePath(path) {
		images = append(images, path)
	}
	return nil
}

func processImages() bktree.BKTree {
	var tree bktree.BKTree
	bar := pb.StartNew(len(images))
	for _, path := range images {
		bar.Increment()
		tree.Add(image{path, hashImage(path)})
	}
	bar.Finish()
	println("Images hashed and BK-tree created")
	return tree
}

func aggregateImages(imagesPath string) {
	fmt.Printf("Starting to cluster your images from %s\n", imagesPath)
	filepath.Walk(imagesPath, traverseDirectory)
	fmt.Printf("Selected %d images\n", len(images))
}

func saveClusters(threshold int, cluster_path string, tree bktree.BKTree) {
	println("Creating clusters")
	clusters := createClusters(tree, images, threshold)
	println("Clusters created")
	createDirectories(clusters, cluster_path)
}

func main() {
	imagesPtr := flag.String("imagesPath", "", "String with the path where all the images are located")
	thresholdPtr := flag.Int("threshold", 10, "Threshold to set similar distances")
	flag.Parse()

	if _, err := os.Stat(*imagesPtr); os.IsNotExist(err) {
		log.Fatal("Given image path does not exist!\n")
	}
	cluster_path := fmt.Sprintf("%s/clusters", *imagesPtr)
	if _, err := os.Stat(cluster_path); os.IsNotExist(err) {
		err := os.Mkdir(cluster_path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	aggregateImages(*imagesPtr)
	tree := processImages()
	saveClusters(*thresholdPtr, cluster_path, tree)
	fmt.Println("Done")
}
