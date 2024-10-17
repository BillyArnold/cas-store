package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "filestore"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:]) //[:] converts bytes to slice

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.filename)
}

type StoreOpts struct {
	// Root is the filder name of the root directory
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	if err != nil {
		return !errors.Is(err, fs.ErrNotExist)
	}
	return true
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformFunc(key)
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FirstPathName())

	defer func() {
		log.Printf("deleted [%s] from disk", pathkey.filename)
	}()

	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(pathAndFilenameWithRoot)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
