package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	HASH_PRINT_LEN = 7
)

type TreeEntry struct {
	Name     string
	Size     int64
	Children []TreeEntry
	State    string
}

func BuildTree(path string) (TreeEntry, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return TreeEntry{}, err
	}

	if !stat.IsDir() {
		return TreeEntry{
			Name: stat.Name(),
			Size: stat.Size(),
		}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return TreeEntry{}, err
	}

	var children []TreeEntry
	var size int64
	for _, entry := range entries {
		c, err := BuildTree(filepath.Join(path, entry.Name()))
		if err != nil {
			return TreeEntry{}, err
		}
		children = append(children, c)
		size += c.Size
	}

	return TreeEntry{
		Name:     stat.Name(),
		Size:     size,
		Children: children,
	}, nil
}

func PrintTree(entry TreeEntry, level int) {
	indent := strings.Repeat(" ", level)
	fmt.Printf("%s%s\n", indent, entry.Name)

	for _, child := range entry.Children {
		PrintTree(child, level+4)
	}
}

//func Compare(source, target string) (bool, TreeEntry, error) {
//	if _, err := os.Stat(source); err != nil {
//		return false, TreeEntry{}, err
//	}
//	if _, err := os.Stat(target); err != nil {
//		return false, TreeEntry{}, err
//	}
//	return false, TreeEntry{}, nil
//}

// ----------------------------------------------------------------------------

func main() {
	//tree, _ := BuildTree(os.Args[1])
	//PrintTree(tree, 0)
	//hash, err := getFilePseudoHash(os.Args[1])
	fmt.Println(CompareFiles(os.Args[1], os.Args[2]))
}
