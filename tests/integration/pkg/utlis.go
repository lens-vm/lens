package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase[TSource any, TResult any] struct {
	// The LensFile contents to test as a string
	LensFile string

	// The set of input values to feed into the lens pipeline
	Input []TSource

	// The set of values expected as output from the lens pipeline
	ExpectedOutput []TResult
}

var hostExecutablePaths = []string{
	getPathRelativeToProjectRoot(
		"/tests/hosts/host-go/build/host-go",
	),
}

func getPathRelativeToProjectRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(path.Dir(filename))))
	return path.Join(root, relativePath)
}

func executeTest[TSource any, TResult any](t *testing.T, testCase TestCase[TSource, TResult]) {
	tempDir := t.TempDir()
	lensFilePath := path.Join(tempDir, "lenseFile.json")
	err := os.WriteFile(lensFilePath, []byte(testCase.LensFile), 0700)
	if err != nil {
		t.Fatal(err)
	}

	inputBytes, err := json.Marshal(testCase.Input)
	if err != nil {
		t.Fatal(err)
	}
	inputJson := string(inputBytes)

	for _, hostPath := range hostExecutablePaths {
		pipeLineCommand := exec.Command(hostPath, lensFilePath)
		pipeLineCommand.Stderr = os.Stderr

		stdin, err := pipeLineCommand.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		var stdout bytes.Buffer
		pipeLineCommand.Stdout = &stdout

		err = pipeLineCommand.Start()
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.WriteString(stdin, inputJson)
		if err != nil {
			t.Fatal(err)
		}

		err = stdin.Close()
		if err != nil {
			t.Fatal(err)
		}

		err = pipeLineCommand.Wait()
		if err != nil {
			t.Fatal(err)
		}

		outputBytes := stdout.Bytes()

		var output []TResult
		err = json.Unmarshal(outputBytes, &output)
		if err != nil {
			t.Fatal(err)
		}

		// We could just assert on the string/byte array, but this gives us clearer errors
		assert.Equal(t, testCase.ExpectedOutput, output)
	}
}
