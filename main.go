package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	_"strconv"
)

/*
	check if an array IndexOf an element:
		in case the element is found returns the index
		in case the element is not found returns -1
*/

func IndexOf[T int|string]( arr []T, val T  ) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}


type Queue []int //Id character

func (q Queue) Add( i int) {
	q = append(q, i)
}

func (q Queue) Pull() {
	
}

type Class struct{
	Name string `json:Name`
} 

type Character struct {
	Id uint `json:"Id"`
	Name string `json:"Name"`
	MaxHp uint `json:"MaxHp"`
	Hp int `json:"Hp"`
	Lvl uint `json:"Lvl"`
	Class int `json:"Class"`
	Friendly bool
}

type Move struct {
	name string
	allowed []int
	desc string
	move func(caster *Character, chs *[]Character) (error)
}

var classes []Class

var moves []Move

func init(){
	ReadClass("files/classes.json")

	moves = []Move{
		Move{
			name: "heal",
			allowed: []int{ classNameToId("Mage")},
			desc: "heals the caster",
			move: func (caster *Character, chs *[]Character) error { 
				(*caster).Hp += 10*int((*caster).Lvl)
				return nil
			},
		},
		Move{
			name: "attack",
			allowed: []int{classNameToId("Mage"),classNameToId("Ranger"),classNameToId("Warrior"),classNameToId("Rogue")},
			desc: "use your weapon to attack one enemy",
			move: func (caster *Character, chs *[]Character) error {
				fmt.Println(*chs)
s
				/* PROOF OF CONCEPT, A REAL API IS NEEDED */
				var prompt string
				for i, v := range *chs {
					// fmt.Println(v.Friendly,caster.Friendly)
					if v.Friendly != caster.Friendly {
						prompt += fmt.Sprintf("\t%d : %s %d\n", i, v.Name, idToClass(v.Class))
					}
				}
				attacked := GetUserInput("who do you want to attack?\n" + prompt)
				/* ------------------------------------- */

				/*if len(*chs) != 1 {
					return fmt.Errorf("%v can only attack one character (DEBUG: attacker %v attack array %v)", (*caster).Id, *caster, *chs )
				}*/
				(*chs)[attacked].Hp -= 10*int((*caster).Lvl)
				return nil
			},
		},
		Move{
			name: "fireball",
			allowed: []int{classNameToId("Mage")},
			desc: "the mage casts a huge fireball, hitting all the enemies",
			move: func (caster *Character, chs *[]Character) error {				
				if len(*chs) < 1 {
					return fmt.Errorf("%v missing enemies characters (DEBUG: attacker %v attack array %v)", (*caster).Id, *caster, chs )
				}
				for i := range *chs{
					if ((*chs)[i].Friendly != caster.Friendly){
						(*chs)[i].Hp -= 10*int((*caster).Lvl)
					}
				}
				return nil
			},
		},
	}
}

func main(){

	characters := []Character{
		Character{Id:0, Name: "pippo",Lvl:2, MaxHp:20, Hp:20, Class:classNameToId("Mage"), Friendly: true},
		Character{Id:1, Name: "taver", Lvl:1, MaxHp:40, Hp:40, Class:classNameToId("Warrior"), Friendly: true},
		Character{Id:2, Name: "mario", Lvl:1, MaxHp:15,  Hp:15, Class:classNameToId("Mage")},
		Character{Id:3, Name: "coca", Lvl:1, MaxHp:20, Hp:20, Class:classNameToId("Rogue")},
		Character{Id:4, Name: "nello", Lvl:1, MaxHp:20, Hp:20, Class:classNameToId("Warrior")},
	}


	for _, v := range characters {
		intiative <- v.Id 
	}
	
	for i:= 0; true; i++{

		IsDead(&characters)

		fmt.Println(" --------- DEBUG --------- ")
		fmt.Println(characters)
		fmt.Println(" --------- DFINE --------- ")
		
		// char := &( characters[i%len(characters)] )

		// sfonnato

		fmt.Println(*char,"Hp del pg: ", (*char).Hp)

		mv := GetUserInput( PrintMoves((*char).Class, moves) + "\n -1 to pass turn" )

		if mv == -1 {
			continue
		}

		// if ((*char).Hp <= 0) {
		// 	fmt.Println("character is death\n")
		// 	characters = characters[]

		// 	continue
		// }
	
		if err := action(moves[mv], char, &characters) ; err != nil{
			fmt.Println(err)
		}
	}

	fmt.Printf( "DEBUG \n %v \n", characters)
}

func IsDead (char *[]Character) {
	var tmp []Character
	for _, v := range *char {
		if v.Hp <= 0 {
			fmt.Println(v, "is dead")
			//*char = append((*char)[:i],(*char)[i+1:]...)
			//fmt.Println(*char)
			continue
		}
		tmp = append(tmp, v)
	
	}

	*char = tmp 
}

func PrintMoves( class int, movest []Move ) (ret string) {
	for i, v := range movest {
		for _, v1 := range v.allowed {
			if class == v1 {
				ret += fmt.Sprintf( "\t %v - %v %v\n", i, v.name, v.desc )
				break
			}
		}
	}
	return
}

/* PROOF OF CONCEPT COMMAND API */
func GetUserInput( prompt string ) (ret int) {

	fmt.Println(prompt)

	fmt.Scan(&ret)
	return 
}
/* ---------------------------- */


func idToClass( i int ) string {
	return classes[i].Name
}

func classNameToId( name string ) int  {
	for i, v := range classes {
		if v.Name == name {
			return i
		}
	}

	return 0
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

func action(move Move, user *Character, targets *[]Character) error {
	if IndexOf(move.allowed, (*user).Class) == -1 {
		return fmt.Errorf("%v is not allowed to use %v", idToClass( (*user).Class ), move.name)
	}
	if err := move.move(user, targets); err != nil {
		fmt.Println(err.Error)
	}
	return nil
}