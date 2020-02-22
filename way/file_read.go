package way

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/thanosKontos/gravelmap"
)

type fileRead struct {
	storageDir string
}

func NewWayFileRead(storageDir string) *fileRead {
	return &fileRead{
		storageDir: storageDir,
	}
}

func (fr *fileRead) Read(ways []gravelmap.Way) ([]string, error) {
	esFl, err := os.Open(fmt.Sprintf("%s/%s", fr.storageDir, edgeStartFilename))
	defer esFl.Close()
	if err != nil {
		return []string{}, err
	}

	plLkFl, err := os.Create(fmt.Sprintf("%s/%s", fr.storageDir, edgeToFilename))
	defer plLkFl.Close()
	if err != nil {
		return []string{}, err
	}

	var edgeFrom int32 = 86123
	var edgeTo int32 = 135138

	nodeStart, err := fr.readEdgeStartFile(esFl, edgeFrom)
	fmt.Println("xxx", nodeStart)


	fr.readEdgeToFile(plLkFl, *nodeStart, edgeTo)



	//for i := 0; int32(i) < nodeStart.ConnectionsCnt; i++ {
	//}

	var polylines []string
	//for _, way := range ways {
	//
	//
	//
	//
	//
	//}


	fmt.Println(ways)

	return polylines, nil
}

func (fr *fileRead) readEdgeToFile(f *os.File, edgeStart edgeStartRecord, edgeToId int32) error {

	//var polygonLookups []int32

	readOffset := edgeStart.NodeToOffset
	f.Seek(readOffset, 0)

	for i := 0; int32(i) < edgeStart.ConnectionsCnt; i++ {
		var storedEdgeTo int32

		data := readNextBytes(f, 4)
		buffer := bytes.NewBuffer(data)
		binary.Read(buffer, binary.BigEndian, &storedEdgeTo)

		fmt.Println("found edge to", storedEdgeTo)

		if storedEdgeTo == edgeToId {

		}

		readOffset = readOffset + 3*4
	}


	fmt.Println("peos")
	os.Exit(0)

	return nil
}

func (fr *fileRead) readEdgeStartFile(f *os.File, edgeStartId int32) (*edgeStartRecord, error) {
	edgeStart := edgeStartRecord{}

	f.Seek(int64(edgeStartId*edgeStartRecordSize), 0)

	data := readNextBytes(f, edgeStartRecordSize)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &edgeStart)
	if err != nil {
		return nil, err
	}

	return &edgeStart, nil
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}


