package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// based on https://gist.github.com/wfjsw/edabe83e3057fbb50c72ca1e5de1d3a5

// MTPParam is a parameter of MTPType
type MTPParam struct {
	Name string
	Type string
}

// MTPType represents a type of the MTProto protocol
type MTPType struct {
	Name       string
	ID         uint32
	Params     []*MTPParam
	ReturnType string // TODO maybe another MTPType?
}

func objectFrom(line string) (*MTPType, error) {
	ignoreLine := line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "---")

	if ignoreLine {
		return nil, nil
	}

	line = strings.Replace(line, ";", "", 1)
	parts := strings.Split(line, " ")

	// extract name and id
	nameWithID := strings.Split(parts[0], "#")
	name := nameWithID[0]
	id := int64(0)
	if len(nameWithID) == 2 {
		id, _ = strconv.ParseInt(nameWithID[1], 16, 64)
	}

	// extract return type
	returnType := parts[len(parts)-1]

	// extract parameters
	params := []*MTPParam{}
	paramsRaw := parts[1 : len(parts)-2]
	if !(strings.HasPrefix(line, "vector")) {
		for _, param := range paramsRaw {
			if strings.HasPrefix(param, "{") {
				continue
			}
			nameWithType := strings.Split(param, ":")
			params = append(params, &MTPParam{
				Name: nameWithType[0],
				Type: nameWithType[1],
			})
		}
	}

	return &MTPType{
		Name:       name,
		ID:         uint32(id),
		ReturnType: returnType,
		Params:     params,
	}, nil
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		return
	}

	lines := bytes.Split(input, []byte("\n"))
	var missedLines []string

	for _, lineBytes := range lines {
		line := string(lineBytes)
		mtpType, err := objectFrom(line)

		if err != nil {
			missedLines = append(missedLines, line)
		} else if mtpType != nil {
			fmt.Printf("id: %x name: %v returnType: %v params: %v\n", mtpType.ID, mtpType.Name, mtpType.ReturnType, mtpType.Params)
		}
	}

	fmt.Println("Missed lines:")
	for _, line := range missedLines {
		fmt.Println(string(line))
	}
}
