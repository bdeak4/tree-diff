package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const HASH_PRINT_LEN = 7

type TreeEntry struct {
	Name     string
	Hash     string
	Size     int64
	Children []TreeEntry
}

func BuildTree(path string) (TreeEntry, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return TreeEntry{}, err
	}

	if !stat.IsDir() {
		hash, err := ComputeHash(path)
		if err != nil {
			return TreeEntry{}, err
		}
		return TreeEntry{
			Name: stat.Name(),
			Size: stat.Size(),
			Hash: hash,
		}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return TreeEntry{}, err
	}

	var children []TreeEntry
	var size int64
	h := sha256.New()
	for _, entry := range entries {
		c, err := BuildTree(filepath.Join(path, entry.Name()))
		if err != nil {
			return TreeEntry{}, err
		}
		children = append(children, c)
		size += c.Size
		h.Write([]byte(c.Hash))
	}

	return TreeEntry{
		Name:     stat.Name(),
		Hash:     fmt.Sprintf("%x", h.Sum(nil)),
		Size:     size,
		Children: children,
	}, nil
}

func PrintTree(entry TreeEntry, level int) {
	indent := strings.Repeat(" ", level)
	fmt.Printf("%s%s %s\n", indent, entry.Hash[:HASH_PRINT_LEN], entry.Name)

	for _, child := range entry.Children {
		PrintTree(child, level+HASH_PRINT_LEN+1)
	}
}

func ComputeHash(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
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

func main() {
	tree, _ := BuildTree(os.Args[1])
	PrintTree(tree, 0)
}
