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

type TreeNode struct {
	Name     string
	Hash     string
	Size     int64
	Children []TreeNode
}

//func Compare(src, dst string) (bool, []string, error) {
//	var diffPaths []string
//
//	if _, err := os.Stat(src); err != nil {
//		return false, diff, err
//	}
//	if _, err := os.Stat(dst); err != nil {
//		return false, diff, err
//	}
//
//	entries, err := os.ReadDir(path)
//}

func BuildTree(path string) (TreeNode, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return TreeNode{}, err
	}

	if !stat.IsDir() {
		hash, err := ComputeHash(path)
		if err != nil {
			return TreeNode{}, err
		}
		return TreeNode{
			Name: stat.Name(),
			Size: stat.Size(),
			Hash: hash,
		}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return TreeNode{}, err
	}

	var children []TreeNode
	var size int64
	h := sha256.New()
	for _, entry := range entries {
		c, err := BuildTree(filepath.Join(path, entry.Name()))
		if err != nil {
			return TreeNode{}, err
		}
		children = append(children, c)
		size += c.Size
		h.Write([]byte(c.Hash))
	}

	return TreeNode{
		Name:     stat.Name(),
		Hash:     fmt.Sprintf("%x", h.Sum(nil)),
		Size:     size,
		Children: children,
	}, nil
}

func PrintTree(entry TreeNode, level int) {
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

func main() {
	tree, _ := BuildTree(os.Args[1])
	PrintTree(tree, 0)
}
