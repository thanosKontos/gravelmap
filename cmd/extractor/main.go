package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"

	//"github.com/golang/protobuf/proto"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "wrong arguments number (expected 1)\n")
		os.Exit(1)
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file\n")
		os.Exit(1)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	fmt.Println(d)
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file\n", err)
		os.Exit(1)
	}

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error decoding file\n", err)
			os.Exit(1)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				//if v.ID == 21487272 || v.ID == 26952503 || v.ID == 4606122609 || v.ID == 2959861566 || v.ID == 270654946 || v.ID == 3949514556 || v.ID == 1312948344 || v.ID == 2959897160 || v.ID == 5955191732 || v.ID == 5955191727 || v.ID == 26952504 {
				//	fmt.Println("Node!", v)
				//}
				//
				//fmt.Println("Node!", v)
				nc++
			case *osmpbf.Way:
				//if v.ID == 4475559 {
				//	fmt.Println("Way!", v)
				//}
				if val, ok := v.Tags["surface"]; ok && val != "asphalt" && val != "cobblestone" {
					fmt.Println("Way!", v)
				}
				//fmt.Println("Way!", v)
				wc++
			case *osmpbf.Relation:
				//fmt.Println("Relation!", v)
				rc++
			default:
				fmt.Fprintf(os.Stderr, "unknown type %T\n", v)
			}
		}
	}

	fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)

	//elliot := &Person{
	//	Name: "Elliot",
	//	Age:  24,
	//}
	//
	//data, err := proto.Marshal(elliot)
	//if err != nil {
	//	fmt.Println("marshalling error")
	//}
	//
	//// printing out our raw protobuf object
	//fmt.Println(data)
	//
	//// let's go the other way and unmarshal
	//// our byte array into an object we can modify
	//// and use
	//newElliot := &Person{}
	//err = proto.Unmarshal(data, newElliot)
	//if err != nil {
	//	fmt.Println("unmarshalling error")
	//}
	//
	//// print out our `newElliot` object
	//// for good measure
	//fmt.Println(newElliot.Age)
	//fmt.Println(newElliot.Name)
}