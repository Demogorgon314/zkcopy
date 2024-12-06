package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

func copyNode(zkSrc *zk.Conn, zkDst *zk.Conn, srcPath string, dstPath string) {
	data, stat, err := zkSrc.Get(srcPath)
	if err != nil {
		log.Fatalf("Failed to get data from source path %s: %v", srcPath, err)
	}

	fmt.Printf("Copying %s ...\n", dstPath)

	// Ensure destination path exists
	parts := strings.Split(dstPath, "/")
	currentPath := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		currentPath += "/" + part
		exists, _, err := zkDst.Exists(currentPath)
		if err != nil {
			log.Fatalf("Failed to check existence of path %s: %v", currentPath, err)
		}
		if !exists {
			_, err := zkDst.Create(currentPath, nil, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				log.Fatalf("Failed to create path %s: %v", currentPath, err)
			}
		}
	}

	// Set data to the destination path
	_, err = zkDst.Set(dstPath, data, -1)
	if err != nil {
		log.Fatalf("Failed to set data for %s: %v", dstPath, err)
	}

	if stat.NumChildren == 0 {
		return
	}

	children, _, err := zkSrc.Children(srcPath)
	if err != nil {
		log.Fatalf("Failed to get children of %s: %v", srcPath, err)
	}

	for _, child := range children {
		srcChildPath := srcPath + "/" + child
		dstChildPath := dstPath + "/" + child
		copyNode(zkSrc, zkDst, srcChildPath, dstChildPath)
	}
}

func main() {
	var srcZk, dstZk, srcPath, dstPath string
	var deleteBeforeCopy bool

	flag.StringVar(&srcZk, "source-zk", "", "Source Zookeeper connection string")
	flag.StringVar(&dstZk, "destination-zk", "", "Destination Zookeeper connection string")
	flag.StringVar(&srcPath, "source-path", "", "Source Zookeeper path")
	flag.StringVar(&dstPath, "destination-path", "", "Destination Zookeeper path")
	flag.BoolVar(&deleteBeforeCopy, "update", false, "Update existing nodes without deletion")
	flag.Parse()

	if srcZk == "" || dstZk == "" || srcPath == "" || dstPath == "" {
		log.Fatalf("ERROR: Missing required arguments")
	}

	if !strings.HasPrefix(srcPath, "/") || !strings.HasPrefix(dstPath, "/") {
		log.Fatalf("ERROR: Paths must start with '/'")
	}

	if len(srcPath) > 1 && strings.HasSuffix(srcPath, "/") {
		log.Fatalf("ERROR: Source path must not end with '/'")
	}

	connSrc, _, err := zk.Connect(strings.Split(srcZk, ","), 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to source Zookeeper: %v", err)
	}
	defer connSrc.Close()

	connDst, _, err := zk.Connect(strings.Split(dstZk, ","), 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to destination Zookeeper: %v", err)
	}
	defer connDst.Close()

	exists, _, err := connSrc.Exists(srcPath)
	if err != nil || !exists {
		log.Fatalf("ERROR: Source path \"%s\" does not exist.", srcPath)
	}

	if deleteBeforeCopy {
		fmt.Printf("Deleting %s ...\n", dstPath)
		err := connDst.Delete(dstPath, -1)
		if err != nil && err != zk.ErrNoNode {
			log.Fatalf("Failed to delete destination path %s: %v", dstPath, err)
		}
	}

	copyNode(connSrc, connDst, srcPath, dstPath)
	fmt.Println("Done.")
}
