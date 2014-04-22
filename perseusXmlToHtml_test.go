package main

import (
        "fmt"
	"testing"
        "strings"
)

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
