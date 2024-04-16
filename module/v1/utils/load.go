package utils

import (
	"Predator/module/v1/ymlpoc/structs"
	"embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

func LoadMultiPoc(Pocs embed.FS, pocname string) []*structs.Poc {
	var pocs []*structs.Poc
	for _, f := range SelectPoc(Pocs, pocname) {
		if p, err := LoadPoc(f, Pocs); err == nil {
			pocs = append(pocs, p)
		}
	}
	return pocs
}

func LoadPoc(fileName string, Pocs embed.FS) (*structs.Poc, error) {
	p := &structs.Poc{}
	yamlFile, err := Pocs.ReadFile("files/" + fileName)

	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		return nil, err
	}
	return p, err
}

func SelectPoc(Pocs embed.FS, pocname string) []string {
	entries, err := Pocs.ReadDir("files")
	if err != nil {
		fmt.Println(err)
	}
	var foundFiles []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), pocname+"-") {
			foundFiles = append(foundFiles, entry.Name())
		}
	}
	return foundFiles
}

func LoadPocbyPath(fileName string) (*structs.Poc, error) {
	p := &structs.Poc{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("[-] load pocs %s error3: %v\n", fileName, err)
		return nil, err
	}
	err = yaml.Unmarshal(data, p)
	if err != nil {
		fmt.Printf("[-] load pocs %s error4: %v\n", fileName, err)
		return nil, err
	}
	return p, err
}
