package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
)

var Classes []Class

type Class struct{
	Id uint `json:Id`
	Name string `json:Name`
}
func main(){
	ReadClass("classes.json")
}

func ReadClass(FileName string){
	content,err :=  ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		return
	}

	err = json.Unmarshal(content, &Classes)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
		return
	}

	// for _, v := range Classes{
	// 	log.Println(v.Id)
	// 	log.Println(v.Name)
	// }
	
}