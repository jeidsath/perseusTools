package main

import (
	"fmt"
	"github.com/jeidsath/unigreek"

	_ "github.com/jeidsath/perseusTools/parse"
)

func main() {
	out, _ := unigreek.Convert("alfa")
	fmt.Println(out)
}
