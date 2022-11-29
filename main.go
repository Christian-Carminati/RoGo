package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	_ "strconv"
)

/*
check if an array IndexOf an element:

	in case the element is found returns the index
	in case the element is not found returns -1
*/

func IndexOf[T any, V any](arr []T, val V, compareFunc func(c1 T, c2 V) bool) int {
	for i, v := range arr {
		if compareFunc(v, val) {
			return i
		}
	}
	return -1
}

func DmgTypeId(name string) int {
	return IndexOf(damageTypes, name, func(v1 DamageType, v2 string) bool {
		return v1.Name == v2
	})
}

func WeaponId(name string) int {
	return IndexOf(weapons, name, func(v1 Weapon, v2 string) bool {
		return v1.Name == v2
	})
}

func ArmorId(name string) int {
	return IndexOf(armors, name, func(v1 Armor, v2 string) bool {
		return v1.Name == v2
	})
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
func raceNameToId(name string) int {
	for i, v := range races {
		if v.Name == name {
			return i
		}
	}
	return 0
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

func prettyPrintStruct[T any](val *[]T) {
	for _, v := range *val {
		s, _ := json.Marshal(v)
		fmt.Println(string(s))
		fmt.Println()
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
	Name        string          `json:"Name"`
	Resistences map[int]float64 `json:"Resistences"`
}

type Race struct {
	Name        string          `json:"Name"`
	Resistences map[int]float64 `json:"Resistences"`
}

type Character struct {
	Id uint `json:"Id"`
	//identify parameters
	Name string `json:"Name"`
	//attributes
	MaxHp uint `json:"MaxHp"`
	Hp    int  `json:"Hp"`
	//flag & status
	Incap    int         `json:"Incap"`
	Status   map[int]int `json:"Status"`
	Focus    bool        `json:"Focus"`
	Friendly bool        `json:"friendly"`
	//skill
	Lvl   uint `json:"Lvl"`
	Class int  `json:"Class"`
	Race  int  `json:"Race"`
	Init  int  `json:"Init"`
	//weapon and armor
	Armor       int             `json:"Armor"`
	Weapon      int             `json:"Weapon"`
	Resistences map[int]float64 // es: 9 5.0 = aggiunge +5 in resistenza nel tipo di danno 9
}

type Move struct {
	Name    string                                                        `json:"Name"`
	Allowed []int                                                         `json:"Allowed"`
	Desc    string                                                        `json:"Desc"`
	Move    func(caster *Character, chs *[]Character, queue *Queue) error `json:"Move"`
}

type StatusEffect struct {
	name      string
	desc      string
	effect    func(key int, caster *Character, chs *[]Character, queue *Queue) error
	endEffect func(key int, caster *Character, chs *[]Character, queue *Queue)
}

type DamageType struct {
	Name string `json:"Name"`
}

type Weapon struct {
	Name       string `json:"Name"`
	DamageType []int
	Damage     int `json:"Damage"`
}

type Armor struct {
	Name string `json:"Name"`
	// hp might not be actually used
	Hp          int             `json:"Hp"`
	Resistences map[int]float64 //dictionary maps a damage resistence and his percentual(1, 0.0)
}

var races []Race
var classes []Class
var damageTypes []DamageType
var moves []Move
var statusEffects []StatusEffect
var armors []Armor
var weapons []Weapon

func init() {

	if err := loadJson("files/classes.json", &classes); err != nil {
		fmt.Printf("Error reading classes %e", err)
	}
	if err := loadJson("files/damage_types.json", &damageTypes); err != nil {
		fmt.Printf("Error reading damage types %e", err)
	}
	if err := loadJson("files/races.json", &races); err != nil {
		fmt.Printf("Error reading damage types %e", err)
	}
	//------------------------------------ FLAGS ------------------------------------------------
	var flag_save = flag.Bool("save", false, "saving json on file")
	flag.Parse()
	if *flag_save {
		serializer()
	}
	//------------------------------------ ENDFLAGS ------------------------------------------------

	if err := loadJson("files/weapons.json", &weapons); err != nil {
		fmt.Printf("Error reading damage types %e", err)
	}
	if err := loadJson("files/armors.json", &armors); err != nil {
		fmt.Printf("Error reading damage types %e", err)
	}

	loadFuncs()
}

func main() {

	characters := []Character{
		{Id: 0, Name: "pippo", Lvl: 2, MaxHp: 20, Hp: 20, Init: 1, Incap: 40, Status: make(map[int]int), Resistences: map[int]float64{9: 5.0, 11: 3.0}, Class: classNameToId("Mage"), Race: raceNameToId("Dwarf"), Weapon: WeaponId("Longsword"), Armor: ArmorId("Old Rusty Chainmail"), Friendly: true},
		{Id: 1, Name: "taver", Lvl: 1, MaxHp: 40, Hp: 40, Init: 6, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior"), Race: raceNameToId("Dwarf"), Weapon: WeaponId("Iron Mace"), Armor: ArmorId("Damaged Plate Armor"), Friendly: true},
		{Id: 2, Name: "mario", Lvl: 1, MaxHp: 15, Hp: 15, Init: 5, Incap: 40, Status: make(map[int]int), Class: classNameToId("Mage"), Race: raceNameToId("Dwarf"), Weapon: WeaponId("Spear"), Armor: ArmorId("Old Rusty Chainmail")},
		{Id: 3, Name: "cocaa", Lvl: 1, MaxHp: 20, Hp: 20, Init: 2, Incap: 30, Status: make(map[int]int), Class: classNameToId("Rogue"), Race: raceNameToId("Dwarf"), Weapon: WeaponId("Crossbow"), Armor: ArmorId("Damaged Plate Armor")},
		{Id: 4, Name: "nello", Lvl: 1, MaxHp: 20, Hp: 20, Init: 4, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior"), Race: raceNameToId("Dwarf"), Weapon: WeaponId("Spear"), Armor: ArmorId("Old Rusty Chainmail")},
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
		prettyPrintStruct(&characters)
		fmt.Println(" --------- DFINE --------- ")

		charIndex, ok := roundQueue.Pull()
		if !ok {
			fmt.Println("something went wrong while pulling new char")
			return
		}

		char := &(characters[charIndex])

		for key := range (*char).Status {
			statusEffects[key].effect(key, char, &characters, roundQueue)
		}

		roundQueue.Add(charIndex)


		// temporary array to print characters
		temporary := &[]Character{ *(char) }
		for _, v := range queue[:len(queue)-1] {
			*(temporary) = append(*(temporary), characters[v])
		}
		printCharacters(temporary)
		temporary = nil
		// stop printing characters

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

		if err := action(moves[mv], char, &characters, roundQueue); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("DEBUG \n %v \n", characters)
}

// using character as argument instead of armor
func calculateDamageProtection(weaponDamageTypes *[]int, chs *Character) float64 {
	//class races armors

	perc := 100.0
	for protType, protVal := range armors[(*chs).Armor].Resistences {
		for _, dmgType := range *weaponDamageTypes {
			if protType == dmgType {
				perc -= (perc / float64(len(*weaponDamageTypes))) * protVal
				break
			}
		}
	}
	for protType, protVal := range classes[(*chs).Class].Resistences {
		for _, dmgType := range *weaponDamageTypes {
			if protType == dmgType {
				perc -= (perc / float64(len(*weaponDamageTypes))) * protVal
				break
			}
		}
	}
	for protType, protVal := range races[(*chs).Class].Resistences {
		for _, dmgType := range *weaponDamageTypes {
			if protType == dmgType {
				perc -= (perc / float64(len(*weaponDamageTypes))) * protVal
				break
			}
		}
	}
	for protType, protVal := range (*chs).Resistences {
		for _, dmgType := range *weaponDamageTypes {
			if protType == dmgType {
				perc -= protVal / float64(len(*weaponDamageTypes))
				break
			}
		}
	}
	return perc / 100.0
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
		for _, v1 := range v.Allowed {
			if class == v1 {
				ret += fmt.Sprintf("\t %d - %s %s\n", i, v.Name, v.Desc)
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
func loadJson[T any](FileName string, inp T) error {
	content, err := os.ReadFile(FileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &inp)
	if err != nil {
		return err
	}

	return nil
}

func action(move Move, user *Character, targets *[]Character, queue *Queue) error {
	if IndexOf(move.Allowed, (*user).Class, func(v1 int, v2 int) bool { return v1 == v2 }) == -1 {
		return fmt.Errorf("%v is not allowed to use %v", idToClass((*user).Class), move.Name)
	}
	if err := move.Move(user, targets, queue); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
