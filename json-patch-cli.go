//json-patch-cli
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/evanphx/json-patch"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func applyPatch(inputs <-chan []byte, outputs chan<- []byte, patch jsonpatch.Patch) {
	in := <-inputs
	out, err := patch.Apply(in)
	check(err)
	outputs <- out
}

func main() {
	var inputPath, patchPath string

	flag.StringVar(&inputPath, "input", "", "[required] path to a JSON file that you want the patch to be be applied to")
	flag.StringVar(&patchPath, "patch", "", "[required] path to a JSON patch file (RFC 6902)")

	flag.Parse()

	input, err := ioutil.ReadFile(inputPath)
	check(err)
	patch, err := ioutil.ReadFile(patchPath)
	check(err)

	patchT, err := jsonpatch.DecodePatch(patch)
	check(err)

	inputs := make(chan []byte, 1)
	outputs := make(chan []byte, 1)

	inputs <- input

	go applyPatch(inputs, outputs, patchT)

	fmt.Println("input:", inputPath)
	fmt.Println("patch:", patchPath)

	fmt.Fprintln(os.Stdout, string(<-outputs))
}
