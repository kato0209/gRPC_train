package main

import (
	"fmt"
	"log"
	"protobuf-lesson/pb"

	"github.com/golang/protobuf/jsonpb"
)

func main() {
	employee := &pb.Employee{
		Id:          1,
		Name:        "Suzuki",
		Email:       "test@example.com",
		Occupation:  pb.Occupation_ENGINEER,
		PhoneNumber: []string{"000-0000-0000", "111-1111-1111"},
		Project:     map[string]*pb.Company_Project{"ProjectX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "my name is Suzuki",
		},
		Birthday: &pb.Date{
			Year:  1990,
			Month: 1,
			Day:   1,
		},
	}

	/*
		binData, err := proto.Marshal(employee)
		if err != nil {
			log.Fatalln("cant serialize", err)
		}

		if err := ioutil.WriteFile("test.bin", binData, 0666); err != nil {
			log.Fatalln("cant write", err)
		}

		in, err := ioutil.ReadFile("test.bin")
		if err != nil {
			log.Fatalln("cant read", err)
		}

		readEmployee := &pb.Employee{}
		err = proto.Unmarshal(in, readEmployee)
		if err != nil {
			log.Fatalln("cant deserialize", err)
		}

		fmt.Println(readEmployee)
	*/

	m := jsonpb.Marshaler{}
	out, err := m.MarshalToString(employee)
	if err != nil {
		log.Fatalln("cant marshal to json", err)
	}
	//fmt.Println(out)

	readEmployee := &pb.Employee{}
	if err := jsonpb.UnmarshalString(out, readEmployee); err != nil {
		log.Fatalln("cant unmarshal from json", err)
	}

	fmt.Println(readEmployee)
}
