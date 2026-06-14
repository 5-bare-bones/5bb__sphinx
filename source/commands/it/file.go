package it

import (
	"fmt"

	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	bolt "go.etcd.io/bbolt"
)

func fileMultiselect(vault *bolt.DB) ([]string, error) {
	files, err := file.ListNames(vault)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		fmt.Println("\nNo files to select")
		return nil, nil
	}

	names, err := multiSelect("Choose files:", files)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func fileMvNames(vault *bolt.DB) ([]string, error) {
	files, err := file.ListNames(vault)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		fmt.Println("\nNo files to select")
		return nil, nil
	}

	// Request src
	src, err := selectOne("Source", "", files)
	if err != nil {
		return nil, err
	}

	// Request dst
	dst, err := input("Destination:", "")
	if err != nil {
		return nil, err
	}

	return []string{src, dst}, nil
}
