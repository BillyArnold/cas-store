package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "thiskeyhere"
	pathkey := CASPathTransformFunc(key)
	expectedOriginalKey := "ec46d6ed1bd314674feaf7989e35e8ff2b27d493"
	expectedPathName := "ec46d/6ed1b/d3146/74fea/f7989/e35e8/ff2b2/7d493"
	if pathkey.Pathname != expectedPathName {
		t.Errorf("have %s want %s", pathkey.Pathname, expectedPathName)
	}
	if pathkey.Original != expectedOriginalKey {
		t.Errorf("have %s want %s", pathkey.Original, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myimagekey", data); err != nil {
		t.Error(err)
	}
}
