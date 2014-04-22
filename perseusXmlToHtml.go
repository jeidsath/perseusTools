package main

import (
	"flag"
	"fmt"
        "io/ioutil"

        "github.com/jeidsath/perseusTools/parse"
)

func ConvertFile(fileName string) ([]string, error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []string{}, err
	}

	root, err := parse.GetTextables(contents)
	if err != nil {
		return []string{}, err
	}

	texts, err := root.Document()
	if err != nil {
		return []string{}, err
	}

	return texts, nil
}

func main() {
	flag.Parse()
        fileName := flag.Arg(0)
	if fileName == "" {
		fmt.Println("No file specified.")
		return
	}

	lines, err := ConvertFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, ss := range lines {
			fmt.Println(ss)
		}
	}
}
