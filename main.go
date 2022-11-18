package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Name        string `json:"Name"`
	Resistences map[int]float64
}

type Race struct {
	Name        string `json:"Name"`
	Resistences map[int]float64
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
	Init  int  `json:"Init"`
	//weapon and armor
	Armor       int `json:"Armor"`
	Weapon      int `json:"Weapon"`
	Resistences map[int]float64
}

type Move struct {
	name    string
	allowed []int
	desc    string
	move    func(caster *Character, chs *[]Character, queue *Queue) error
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
	armors = []Armor{
		{
			Name: "Old Rusty Chainmail",
			Resistences: map[int]float64{
				DmgTypeId("Slashing"): 0.6,
				DmgTypeId("Piercing"): 0.2,
			},
		},
		{
			Name: "Damaged Plate Armor",
			Resistences: map[int]float64{
				DmgTypeId("Piercing"):    0.6,
				DmgTypeId("Bludgeoning"): 0.4,
				DmgTypeId("Slashing"):    0.2,
			},
		},
	}
	weapons = []Weapon{
		{
			Name: "Longsword",
			DamageType: []int{
				DmgTypeId("Slashing"),
			},
			Damage: 8,
		},
		{
			Name: "Spear",
			DamageType: []int{
				DmgTypeId("Piercing"),
			},
			Damage: 4,
		},
		{
			Name: "Iron Mace",
			DamageType: []int{
				DmgTypeId("Bludgeoning"),
			},
			Damage: 15,
		},
		{
			Name: "Crossbow",
			DamageType: []int{
				DmgTypeId("Piercing"),
			},
			Damage: 10,
		},
	}

	moves = []Move{
		{
			name:    "self-heal",
			allowed: []int{classNameToId("Mage")},
			desc:    "heals the caster",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {
				// caster heals himself
				if (*caster).Hp+(10*int((*caster).Lvl)) > int((*caster).MaxHp) {
					(*caster).Hp = int((*caster).MaxHp)
					return nil
				}
				(*caster).Hp += 10 * int((*caster).Lvl)
				return nil
			},
		},
		{
			name:    "attack",
			allowed: []int{classNameToId("Mage"), classNameToId("Ranger"), classNameToId("Warrior"), classNameToId("Rogue")},
			desc:    "use your weapon to attack one enemy",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {
				var Damage = weapons[(*caster).Weapon].Damage
				var DamageType = weapons[(*caster).Weapon].DamageType
				// character uses his melee weapon to attack an enemy

				//fmt.Println(*chs)
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

				//fmt.Println(int(calculateDamageProtection(&weapons[(*caster).Weapon].DamageType,&armors[(*chs)[attacked].Armor]) * float64(weapons[(*caster).Weapon].Damage) * float64((*caster).Lvl)))
				(*chs)[attacked].Hp -= int(calculateDamageProtection(&DamageType, &(*chs)[attacked]) * float64(Damage) * float64((*caster).Lvl))
				return nil
			},
		},
		{
			name:    "fireball",
			allowed: []int{classNameToId("Mage")},
			desc:    "the mage casts a huge fireball, hitting all the enemies",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {

				var Damage = 10
				var DamageType = []int{DmgTypeId("Fire")}

				// fireball deals AOE damage, it also targets the dead

				// not even sure this is needed
				if len(*chs) < 1 {
					return fmt.Errorf("%v missing enemies characters (DEBUG: attacker %v attack array %v)", (*caster).Id, *caster, chs)
				}
				//var DamageType
				for i := range *chs {
					if (*chs)[i].Friendly != caster.Friendly {

						//fmt.Println(int(calculateDamageProtection(&DamageType,&armors[(*chs)[i].Armor])))
						(*chs)[i].Hp -= Damage * int(calculateDamageProtection(&DamageType, &(*chs)[i])) * int((*caster).Lvl)
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
				(*caster).Status[2] = int((*chs)[i].Id)
				(*caster).Focus = true

				return nil
			},
		},
		{
			name:    "poisonus dart",
			allowed: []int{classNameToId("Rogue")},
			desc:    "the attacker launches a poisoned dart, dealing 5 dmg and posioning the subject for 2*Lvl stacks",
			move: func(caster *Character, chs *[]Character, queue *Queue) error {

				var DamageArrow = weapons[(*caster).Weapon].Damage
				var DamageTypeArrow = weapons[(*caster).Weapon].DamageType
				var stack = 5

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

				(*chs)[i].Hp -= int(calculateDamageProtection(&DamageTypeArrow, &(*chs)[i]) * float64(DamageArrow) * float64((*caster).Lvl))

				if (*chs)[i].Status == nil {
					(*chs)[i].Status = make(map[int]int)
				}
				(*chs)[i].Status[0] = stack * int(caster.Lvl)

				return nil
			},
		},
	}

	statusEffects = []StatusEffect{
		{
			name: "poison",
			desc: "the character is poisoned, taking damage every turn",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {

				if (*caster).Status[key] <= 0 {
					statusEffects[key].endEffect(key, caster, chs, queue)
				}

				(*caster).Hp -= (*caster).Status[key]

				(*caster).Status[key]--

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

				//fmt.Println("DEBUG: ", *caster)
				(*caster).Focus = true

				if (*caster).Status[key] <= 0 {
					//(*caster).Friendly = !(*caster).Friendly
					statusEffects[key].endEffect(key, caster, chs, queue)
				}
				(*caster).Status[key]--

				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				(*caster).Friendly = !(*caster).Friendly
				(*caster).Focus = false

				/*fmt.Println("dehhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
				fmt.Println(caster)*/

				for i := range *chs {
					//val, ok := (*chs)[i].Status[2]
					//fmt.Println(val, ok, "|", (*caster).Id )
					if val, ok := (*chs)[i].Status[2]; ok && val == int((*caster).Id) {
						//fmt.Println((*chs)[i])
						statusEffects[2].endEffect(2, &((*chs)[i]), chs, queue)
					}
				}

				/*fmt.Println("dehhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
				fmt.Printf("Porcoddio: %p\n", caster)
				fmt.Println((*caster), "| deleting key: ", key )*/
				delete((*caster).Status, key)
				//fmt.Println((*caster))
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
	if err := WriteJson("files/armors.json", &armors); err != nil {
		fmt.Printf("Error writing armors %e", err)
	}
	if err := WriteJson("files/weapons.json", &weapons); err != nil {
		fmt.Printf("Error writing weapons %e", err)
	}
}

func main() {

	for protType, protVal := range classes[1].Resistences {
		fmt.Println(protType, protVal)
	}

	characters := []Character{
		{Id: 0, Name: "pippo", Lvl: 2, MaxHp: 20, Hp: 20, Init: 1, Incap: 40, Status: make(map[int]int), Class: classNameToId("Mage"), Weapon: WeaponId("Longsword"), Armor: ArmorId("Old Rusty Chainmail"), Friendly: true},
		{Id: 1, Name: "taver", Lvl: 1, MaxHp: 40, Hp: 40, Init: 6, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior"), Weapon: WeaponId("Iron Mace"), Armor: ArmorId("Damaged Plate Armor"), Friendly: true},
		{Id: 2, Name: "mario", Lvl: 1, MaxHp: 15, Hp: 15, Init: 5, Incap: 40, Status: make(map[int]int), Class: classNameToId("Mage"), Weapon: WeaponId("Spear"), Armor: ArmorId("Old Rusty Chainmail")},
		{Id: 3, Name: "cocaa", Lvl: 1, MaxHp: 20, Hp: 20, Init: 2, Incap: 30, Status: make(map[int]int), Class: classNameToId("Rogue"), Weapon: WeaponId("Crossbow"), Armor: ArmorId("Damaged Plate Armor")},
		{Id: 4, Name: "nello", Lvl: 1, MaxHp: 20, Hp: 20, Init: 4, Incap: 10, Status: make(map[int]int), Class: classNameToId("Warrior"), Weapon: WeaponId("Spear"), Armor: ArmorId("Old Rusty Chainmail")},
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
			//val := (*char).Status[key]
			/*if val == 0 {
				statusEffects[key].endEffect(key, char, &characters, roundQueue)
				continue
			}*/
			//fmt.Printf("Porcoddio: %p\n", char)
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

//TODO: missing natural protection given by race or class
/*
	calculates dmg protection given by given armor
*/

// using character as argument instead of armor
func calculateDamageProtection(weaponDamageTypes *[]int, chs *Character) float64 {
	//fmt.Println((*weapon), (*armor),(*armor).Resistences)
	perc := 100.0
	for protType, protVal := range armors[(*chs).Armor].Resistences {
		for _, dmgType := range *weaponDamageTypes {
			if protType == dmgType {
				//fmt.Println(perc,len(*weaponDamageTypes), (*weaponDamageTypes), (*armor).Resistences ,protVal)
				perc -= (100 / float64(len(*weaponDamageTypes))) * protVal
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

	return fmt.Sprintf(" %c  lvl %d | %s | %s | %s | %T ", isAlly, char.Lvl, char.Name, idToClass(char.Class), HpStatus, char.Status)
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
func WriteJson[T any](FileName string, inp T) error {
	file, _ := json.MarshalIndent(inp, "", "\t")

	err := ioutil.WriteFile(FileName, file, 0644)
	return err
}

func loadJson[T any](FileName string, inp T) error {
	content, err := os.ReadFile(FileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &inp)
	if err != nil {
		return err
	}

	// for _, v := range Classes{
	// 	log.Println(v.Id)
	// 	log.Println(v.Name)
	// }
	return nil
}

func action(move Move, user *Character, targets *[]Character, queue *Queue) error {
	if IndexOf(move.allowed, (*user).Class, func(v1 int, v2 int) bool { return v1 == v2 }) == -1 {
		return fmt.Errorf("%v is not allowed to use %v", idToClass((*user).Class), move.name)
	}
	if err := move.move(user, targets, queue); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
