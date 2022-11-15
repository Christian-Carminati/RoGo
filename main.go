package main

import (
	"encoding/json"
	"fmt"
	"os"
	_ "strconv"
)

/*
check if an array IndexOf an element:

	in case the element is found returns the index
	in case the element is not found returns -1
*/

func IndexOf[T int|string](arr []T, val T) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

func bubbleSort[T any](arr *[]T, compare func(c1 T, c2 T) bool) {
	n := len(*arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if compare((*arr)[j], (*arr)[j+1]) {
				(*arr)[j], (*arr)[j+1] = (*arr)[j+1], (*arr)[j]
			}
		}
	}
}

type Queue []int //Id character in fight

func (q *Queue) Add(i int) {
	*q = append(*q, i)
}

func (q *Queue) Pull() (ret int, ok bool) {

	if len(*q) == 0 {
		ret = -1
		ok = false
		return
	}

	ret = (*q)[0]
	ok = true
	*q = (*q)[1:]
	return
}

type HpStatus int

const (
	Alive HpStatus = iota // 0
	Incap
	Dead
	Mutil
)

type Class struct {
	Name string `json:"Name"`
}

type Character struct {
	Id uint `json:"Id"`

	Name string `json:"Name"`

	MaxHp  uint `json:"MaxHp"`
	Hp     int  `json:"Hp"`
	Incap  int  `json:"Incap"`
	Status map[int]int
	Focus bool `json:"Focus"`

	Lvl   uint `json:"Lvl"`
	Class int  `json:"Class"`
	Init  int  `json:"Init"`

	Friendly bool
}

type Move struct {
	name    string
	allowed []int
	desc    string
	move    func(caster *Character, chs *[]Character, queue *Queue) error
}

type StatusEffect struct {
	name   string
	desc   string
	effect func(key int, caster *Character, chs *[]Character, queue *Queue) error
	endEffect func(key int, caster *Character, chs *[]Character, queue *Queue)
}

var classes []Class

var moves []Move
var statusEffects []StatusEffect

func init() {
	ReadClass("files/classes.json")

	moves = []Move{
		{
			name:    "self-heal",
			allowed: []int{classNameToId("Mage")},
			desc:    "heals the caster",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {
				// caster heals himself

				(*caster).Hp += 10 * int((*caster).Lvl)
				return nil
			},
		},
		{
			name:    "attack",
			allowed: []int{classNameToId("Mage"), classNameToId("Ranger"), classNameToId("Warrior"), classNameToId("Rogue")},
			desc:    "use your weapon to attack one enemy",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {
				// character uses his melee weapon to attack an enemy

				fmt.Println(*chs)
				/* PROOF OF CONCEPT, A REAL API IS NEEDED */
				var prompt string
				for i, v := range *chs {
					// fmt.Println(v.Friendly,caster.Friendly)
					if v.Friendly != caster.Friendly && v.Hp > 0-int(v.MaxHp) {
						prompt += fmt.Sprintf("\t%d : %s \n", i, formatChar(v)) //fmt.Sprintf("\t%d : %s %s\n", i, v.Name, idToClass(v.Class))
					}
				}
				attacked := GetUserInput("who do you want to attack?\n" + prompt)
				/* ------------------------------------- */

				/*if len(*chs) != 1 {
					return fmt.Errorf("%v can only attack one character (DEBUG: attacker %v attack array %v)", (*caster).Id, *caster, *chs )
				}*/
				(*chs)[attacked].Hp -= 10 * int((*caster).Lvl)
				return nil
			},
		},
		{
			name:    "fireball",
			allowed: []int{classNameToId("Mage")},
			desc:    "the mage casts a huge fireball, hitting all the enemies",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {
				// fireball deals AOE damage, it also targets the dead

				// not even sure this is needed
				if len(*chs) < 1 {
					return fmt.Errorf("%v missing enemies characters (DEBUG: attacker %v attack array %v)", (*caster).Id, *caster, chs)
				}

				for i := range *chs {
					if (*chs)[i].Friendly != caster.Friendly {
						(*chs)[i].Hp -= 10 * int((*caster).Lvl)
					}
				}
				return nil
			},
		},
		{
			name:    "mind control",
			allowed: []int{classNameToId("Mage")},
			desc:    "the mage controls the mind of the enemy for ⌊lvl/2⌋+1 turns",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {

				if (*caster).Focus == true {
					return fmt.Errorf("caster does not have focus")
				}

				/* PROOF OF CONCEPT, A REAL API IS NEEDED */
				var prompt string
				for i, v := range *chs {
					// fmt.Println(v.Friendly,caster.Friendly)
					_, okMindC := v.Status[1]
					_, okCMind := v.Status[2]
					if v.Friendly != caster.Friendly && v.Hp > 0-int(v.MaxHp) && !okMindC && !okCMind {
						prompt += fmt.Sprintf("\t%d : %s \n", i, formatChar(v)) //fmt.Sprintf("\t%d : %s %s\n", i, v.Name, idToClass(v.Class))
					}
				}
				i := GetUserInput("who do you want to control?\n" + prompt)
				/* ------------------------------------- */

				(*chs)[i].Friendly = !(*chs)[i].Friendly
				(*chs)[i].Focus = true

				if (*chs)[i].Status == nil {
					(*chs)[i].Status = make(map[int]int)
				}
				if (*caster).Status == nil {
					(*caster).Status = make(map[int]int)
				}
				(*chs)[i].Status[1] = int(caster.Lvl/2) + 1
				(*caster).Status[2] = int((*chs)[i].Id)+1
				(*caster).Focus = true

				return nil
			},
		},
		{
			name:    "poisonus dart",
			allowed: []int{classNameToId("Rogue")},
			desc:    "the attacker launches a poisoned dart, dealing 5 dmg and posioning the subject for 2*Lvl turns",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {

				/* PROOF OF CONCEPT, A REAL API IS NEEDED */
				var prompt string
				for i, v := range *chs {
					// fmt.Println(v.Friendly,caster.Friendly)
					if v.Friendly != caster.Friendly && v.Hp > 0-int(v.MaxHp) {
						prompt += fmt.Sprintf("\t%d : %s \n", i, formatChar(v)) //fmt.Sprintf("\t%d : %s %s\n", i, v.Name, idToClass(v.Class))
					}
				}
				i := GetUserInput("who do you want to attack?\n" + prompt)
				/* ------------------------------------- */

				(*chs)[i].Hp -= 5

				if (*chs)[i].Status == nil {
					(*chs)[i].Status = make(map[int]int)
				}
				(*chs)[i].Status[0] = int(caster.Lvl)

				return nil
			},
		},
	}

	statusEffects = []StatusEffect{
		{
			name: "poison",
			desc: "the character is poisoned, taking damage every turn",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {

				(*caster).Hp -= (*caster).Status[key]

				(*caster).Status[key]--

				if ((*caster).Status[key] <= 0) {
					statusEffects[key].endEffect(key, caster, chs, queue)
				}

				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				delete((*caster).Status, key)
			},
		},
		{
			name: "mind control",
			desc: "the character changes factions",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {

				(*caster).Status[key]--

				if (*caster).Status[key] == 0 {
					//(*caster).Friendly = !(*caster).Friendly
					statusEffects[key].endEffect(key, caster, chs, queue)
				}
				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				(*caster).Friendly = !(*caster).Friendly
				(*caster).Focus = false

				fmt.Println("dehhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
				fmt.Println(caster)

				for i := range (*chs){
					fmt.Println((*chs)[i].Status[2])
					if val, ok := (*chs)[i].Status[2]; ok && val == int((*caster).Id) {
						fmt.Println((*chs)[i])
						statusEffects[2].endEffect(key, &((*chs)[i]), chs, queue)
					}
				}

				fmt.Println("dehhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
				delete((*caster).Status, key)
			},
		},
		{
			name: "controlling mind",
			desc: "caster is controlling the mind of another character",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {
				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				(*caster).Focus = false
				delete((*caster).Status, key)
			},
		},
	}
}

func main() {

	characters := []Character{
		{Id: 0, Name: "pippo", Lvl: 2, MaxHp: 20, Hp: 20, Init: 1, Incap: 40, Status: make(map[int]int), Class: classNameToId("Mage"), Friendly: true},
		{Id: 1, Name: "taver", Lvl: 1, MaxHp: 40, Hp: 40, Init: 6, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior"), Friendly: true},
		{Id: 2, Name: "mario", Lvl: 1, MaxHp: 15, Hp: 15, Init: 5, Incap: 40, Status: make(map[int]int), Class: classNameToId("Mage")},
		{Id: 3, Name: "coca", Lvl: 1, MaxHp: 20, Hp: 20, Init: 2, Incap: 30, Status: make(map[int]int), Class: classNameToId("Rogue")},
		{Id: 4, Name: "nello", Lvl: 1, MaxHp: 20, Hp: 20, Init: 4, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior")},
	}

	var queue Queue
	roundQueue := &queue

	bubbleSort(&characters, func(c1 Character, c2 Character) bool {
		return c1.Init < c2.Init
	})

	for i := range characters {
		roundQueue.Add(i)
	}

	for i := 0; true; i++ {

		fmt.Println(" --------- DEBUG --------- ")
		fmt.Println(characters)
		fmt.Println(" --------- DFINE --------- ")

		charIndex, ok := roundQueue.Pull()
		if !ok {
			fmt.Println("something went wrong while pulling new char")
			return
		}

		char := &(characters[charIndex])

		for key := range (*char).Status {
			//val := (*char).Status[key]
			/*if val == 0 {
				statusEffects[key].endEffect(key, char, &characters, roundQueue)
				continue
			}*/
			statusEffects[key].effect(key, char, &characters, roundQueue)
		}

		roundQueue.Add(charIndex)


		fmt.Println()
		fmt.Println(0, ":"+formatChar(*char))
		for i, v := range queue[:len(queue)-1] {
			fmt.Println(i+1, ":"+formatChar(characters[v]))
		}

		fmt.Println()


		// in case the character is dead just skip his turn
		if userHpStatus(*char) > 0 {
			continue
		}

		if FightIsOver(&characters) {
			fmt.Println("The fight is Over", characters)
			break
		}

		// sfonnato

		fmt.Println(*char, "Hp del pg: ", (*char).Hp)

		mv := GetUserInput(PrintMoves((*char).Class, moves) + "\n -1 to pass turn")

		if mv == -1 {
			continue
		}

		// if ((*char).Hp <= 0) {
		// 	fmt.Println("character is death\n")
		// 	characters = characters[]

		// 	continue
		// }

		if err := action(moves[mv], char, &characters, roundQueue); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("DEBUG \n %v \n", characters)
}

func FightIsOver(char *[]Character) bool {

	var faction bool
	var valid bool
	for _, v := range *char {

		if v.Hp > 0 && valid == false {
			faction = v.Friendly
			valid = true
		}

		if valid && v.Friendly != faction && v.Hp > 0 {
			return false
		}
	}

	// clean up character array

	return true
}

func IncapDmg(maxHp uint, incap int) int {
	return int(float64(maxHp) * (float64(incap) / 100))
}

func userHpStatus(char Character) HpStatus {

	switch {
	case char.Hp <= 0-int(char.MaxHp):
		return Mutil
	case char.Hp <= 0:
		return Dead
	case char.Hp <= IncapDmg(char.MaxHp, char.Incap):
		return Incap
	default:
		return Alive
	}
}

func PrintMoves(class int, movest []Move) (ret string) {
	for i, v := range movest {
		for _, v1 := range v.allowed {
			if class == v1 {
				ret += fmt.Sprintf("\t %d - %s %s\n", i, v.name, v.desc)
				break
			}
		}
	}
	return
}

/* PROOF OF CONCEPT COMMAND API */
func GetUserInput(prompt string) (ret int) {

	fmt.Println(prompt)

	fmt.Scan(&ret)
	return
}

/* ---------------------------- */

func formatChar(char Character) string {

	var HpStatus string
	uhs := userHpStatus(char)

	switch uhs {
	case Mutil:
		HpStatus = "MUTIL"
	case Dead:
		HpStatus = "DEAD "
	case Incap:
		HpStatus = "INCAP"
	default:
		if char.Hp == int(char.MaxHp) {
			// char has not been hit
			HpStatus = "NOHIT"
		} else if char.Hp > int(float64(char.MaxHp)*0.66) {
			// char is lightly damaged
			// at this stage mages and other fragile classes are already almost incapacitated
			HpStatus = "DAMGD"
		} else if char.Hp > int(float64(char.MaxHp)*0.33) {
			// char is wounded
			// this mostly applies to tough classes
			// at this stage median classes like the ranger are almost incapacitated
			HpStatus = "WOUND"
		} else {
			// char is at the dead door
			// the character is basically dead
			HpStatus = "DDOOR"
		}
	}

	isAlly := ' '
	if char.Friendly {
		isAlly = '⚝'
	}

	return fmt.Sprintf(" %c  lvl %d | %s | %s | %s | %s ", isAlly, char.Lvl, char.Name, idToClass(char.Class), HpStatus, char.Status)
}

func idToClass(i int) string {
	return classes[i].Name
}

func classNameToId(name string) int {
	for i, v := range classes {
		if v.Name == name {
			return i
		}
	}

	return 0
}

func ReadClass(FileName string) {
	content, err := os.ReadFile(FileName)
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

func action(move Move, user *Character, targets *[]Character, queue *Queue) error {
	if IndexOf(move.allowed, (*user).Class) == -1 {
		return fmt.Errorf("%v is not allowed to use %v", idToClass((*user).Class), move.name)
	}
	if err := move.move(user, targets, queue); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
