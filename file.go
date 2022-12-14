package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

const (
	PSEUDOHASH_CHUNK_SIZE = 256 * 1024
)

type Diff struct {
	oldValue, newValue interface{}
}

func CompareFiles(path1, path2 string) (bool, []Diff, error) {
	size1, err := getFileSize(path1)
	if err != nil {
		return false, []Diff{}, err
	}

	size2, err := getFileSize(path2)
	if err != nil {
		return false, []Diff{}, err
	}

	if size1 != size2 {
		sizediff := Diff{size1, size2}
		return false, []Diff{sizediff}, err
	}

	hashFunc := getFileHash
	if size1 > 2*PSEUDOHASH_CHUNK_SIZE {
		hashFunc = getFilePseudoHash
	}

	hash1, err := hashFunc(path1)
	if err != nil {
		return false, []Diff{}, err
	}

	hash2, err := hashFunc(path2)
	if err != nil {
		return false, []Diff{}, err
	}

	if hash1 != hash2 {
		sizediff := Diff{size1, size2}
		hashdiff := Diff{hash1, hash2}
		return false, []Diff{sizediff, hashdiff}, err
	}

	return true, []Diff{}, nil
}

func getFileSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func getFilePseudoHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.CopyN(h, f, PSEUDOHASH_CHUNK_SIZE); err != nil {
		return "", err
	}

	f.Seek(-PSEUDOHASH_CHUNK_SIZE, io.SeekEnd)
	if _, err := io.CopyN(h, f, PSEUDOHASH_CHUNK_SIZE); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
