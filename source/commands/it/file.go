package it

import (
	"fmt"

	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/AlecAivazis/survey/v2"
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

	namesQs := []*survey.Question{
		{
			Name: "names",
			Prompt: &survey.MultiSelect{
				Message: "Choose files:",
				Options: files,
				VimMode: true, //TODO: remove VimMode
			},
		},
	}

	names := []string{}
	if err := ask(namesQs, &names); err != nil {
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
	qs := selectQs("Source", "", files)
	src := struct{ Name string }{}
	if err := ask(qs, &src); err != nil {
		return nil, err
	}

	// Request dst
	dstQs := &survey.Input{
		Message: "Destination:",
	}
	dst, err := askOne(dstQs)
	if err != nil {
		return nil, err
	}

	return []string{src.Name, dst}, nil
}
