package main

import (
	"fmt"
	"os"

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

	fmt.Print(f.Stat())

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