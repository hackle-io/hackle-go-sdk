package bucketer

import (
	"bufio"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMurMur3Hasher(t *testing.T) {
	sut := murmur3Hasher{}

	test := func(filename string) {
		readFile, err := os.Open("../../../../testdata/" + filename + ".csv")
		//goland:noinspection GoUnhandledErrorResult
		defer readFile.Close()

		if err != nil {
			fmt.Println(err)
		}

		fs := bufio.NewScanner(readFile)
		fs.Split(bufio.ScanLines)

		for fs.Scan() {
			line := fs.Text()
			row := strings.Split(line, ",")
			s := i32(row[1])
			h := i32(row[2])
			actual := sut.Hash(row[0], s)
			assert.Equal(t, h, actual)
		}

	}
	test("murmur_all")
	test("murmur_alphabetic")
	test("murmur_alphanumeric")
	test("murmur_numeric")
	test("murmur_uuid")
}

func i32(value string) int32 {
	number, ok := types.AsNumber(value)
	if !ok {
		panic(value)
	}
	return int32(number)
}
