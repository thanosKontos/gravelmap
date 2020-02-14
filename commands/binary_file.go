package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// binaryFileCommand...
func binaryFileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "binary-test",
		Short: "a bin file test",
		Long:  "a bin file test",
		Run: func(cmd *cobra.Command, args []string) {
			writeFile()
			readFile()
		},
	}
}

func readFile() {
	file, err := os.Open("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var m int64

	//file.Seek(16, 0)
	file.Seek(8000, 0)
	data := readNextBytes(file, 8)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m)
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func writeFile() {
	file, err := os.Create("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var t1 int64 = 1234
	var t2 int64 = 5678
	var t3 int64 = 4321

	file.Seek(8000, 0)

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, t1)
	writeNextBytes(file, bin_buf.Bytes())

	var bin_buf2 bytes.Buffer
	binary.Write(&bin_buf2, binary.BigEndian, t2)
	writeNextBytes(file, bin_buf2.Bytes())

	var bin_buf3 bytes.Buffer
	binary.Write(&bin_buf3, binary.BigEndian, t3)
	writeNextBytes(file, bin_buf3.Bytes())
}
func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}
