// Copyright 2020 Google Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"bufio"
)


func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	//yaml file  to filetextlines
	readFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
 
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
 
	readFile.Close()
 
	// decode
	var v yaml.Node
	if err := yaml.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	// TODO investigate how are multiple docs coming to Node
	if len(v.Content) == 0 {
		panic("no yaml docs found")
	}

	// TODO support multiple documents
	content := v.Content[0]

	colorizeKeys(content, "$root", fileTextLines)


	for _,s := range fileTextLines{
		fmt.Println(s)	
	}
}


func colorizeKeys(node *yaml.Node, path string,fileTextLines []string) {
	var prevKey string
	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}
		if i%2 == 0 && child.Value != "" {
			keyPath := path + "." + child.Value
			prevKey = child.Value
			colorAttribute := colorForKey(keyPath)
			if colorAttribute != color.FgBlack {
				addcolor(fileTextLines,colorAttribute,child)
			}
		}

		subPath := path
		if node.Kind == yaml.MappingNode {
			subPath = path + "." + prevKey
		}
		colorizeKeys(child, subPath, fileTextLines)
	}
}


func addcolor(fileTextLines []string,colorAttribute color.Attribute,child *yaml.Node){
	mkcolor := color.New(colorAttribute, color.Bold).SprintFunc()
	line := child.Line-1
	start := child.Column-1
	end := start+len(child.Value)
	fileTextLines[line] = fileTextLines[line][0:start] + mkcolor(fileTextLines[line][start:end]) + fileTextLines[line][end:] 
}

func colorForKey(path string) color.Attribute {
	redSuffixes := []string{"$root.apiVersion",
		"$root.kind",
		"$root.metadata",
		".spec",
		".containers.name",
		".containers.image"}
	for _, f := range redSuffixes {
		if strings.HasSuffix(path, f) {
			return color.FgRed
		}
	}

	if strings.HasPrefix(path, "$root.metadata") {
		return color.FgYellow
	}

	if strings.HasPrefix(path, "$root.spec") {
		return color.FgBlue
	}
	return color.FgBlack
}
