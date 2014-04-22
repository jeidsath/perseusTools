package main

import (
        "fmt"
	"testing"
        "io/ioutil"
        "strings"
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

func TestInitial(t *testing.T) {
        texts, err := ConvertFile("testData/xen.anab_gk.xml")

        if err != nil {
                fmt.Println(err.Error())
                t.Fail()
        }

        for _, ss := range texts {
                fmt.Println(ss)
        }
}

func TestCodes(t *testing.T) {
        texts, err := ConvertFile("testData/xen.anab_gk.xml")

        if err != nil {
                fmt.Println(err.Error())
                t.Fail()
        }

        for _, ss := range texts {
                if strings.ContainsRune(ss, '&') {
                        fmt.Printf("Contains escape: %s\n", ss)
                        t.Fail()
                }
        }
}
