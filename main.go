package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

/*
	check if an array contains an element:
		in case the element is found returns the index
		in case the element is not found returns -1
*/
func Contains[T int|string]( arr []T, val T  ) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

type Class struct{
	Name string `json:Name`
}

type Character struct {
	Id uint `json:"Id"`
	Name string `json:"Name"`
	Hp int `json:"Hp"`
	Lvl uint `json:"Lvl"`
	Class int `json:"Class"`
}

type Move struct {
	name string
	allowed []int
	desc string
	move func(ch1 *Character, ch2 *Character)
}

var classes []Class

var moves = []Move{
	Move{
		name: "heal",
		allowed: []int{ 1 },
		desc: "heals the caster",
		move: func (ch1 *Character, ch2 *Character) { 
			ch1.Hp += 10*int(ch1.Lvl)
		},
	},
	Move{
		name: "attack",
		allowed: []int{ 1, 2, 3, 4 },
		desc: "use your weapon to attack the enemy",
		move: func (ch1 *Character, ch2 *Character) { 
			ch2.Hp -= 10*int(ch1.Lvl)
		},
	},
}

func main(){

	ReadClass("files/classes.json")

	c1 := Character{Id:1, Lvl:1, Hp:10, Class:1}
	c2 := Character{Id:2, Lvl:1, Hp:15, Class:1}
	
	if err := action(moves[0], &c1, &c2 ); err != nil{
		fmt.Println(err)
	}

	fmt.Println(c1)
	fmt.Println(c2)

	// for{
		
	// }
}

func idToClass( i int ) string {
	return classes[i].Name
}

func ReadClass(FileName string){
	content,err :=  ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		return
	}

	err = json.Unmarshal(content, &classes)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
		return
	}

	// for _, v := range Classes{
	// 	log.Println(v.Id)
	// 	log.Println(v.Name)
	// }
	
}

func action(move Move, user *Character, target *Character) error {
	if Contains(move.allowed, (*user).Class) == -1 {
		return fmt.Errorf("%v is not allowed to use %v", idToClass( (*user).Class ), move.name)
	}
	move.move(user, target)
	return nil
}