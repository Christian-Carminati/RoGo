package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"atomicgo.dev/cursor"
)

var regexMatchSingleChar *regexp.Regexp
var regexMatchMultipleChars *regexp.Regexp
var regexSplitWs *regexp.Regexp

func init() {
	regexMatchSingleChar = regexp.MustCompile(`^\s*(?P<number>\d+)\s*\n$`)
	regexMatchMultipleChars = regexp.MustCompile(`^(?P<args>[\d\s]*)\n$`)
	regexSplitWs = regexp.MustCompile(`\s+`)
}

func SingleSelector[T any, S interface{}](title string, elements *[]T, vars S, filter func(element T, vars S) bool) int {

	fmt.Println(title)
	cont := 0
	mapRetIndex := make(map[int]int)
	for i, v := range *elements {
		if filter(v, vars) {
			fmt.Print("\t", cont, ": ", v, "\n")
			mapRetIndex[cont] = i
			cont++
		}
	}

	var inp int
	var opt string
	for {
		inp = -1
		fmt.Print("\r" + opt + "press a number to select: ")
		r := bufio.NewReader(os.Stdin)
		input, _ := r.ReadString('\n')

		if regexMatchSingleChar.MatchString(input) {
			match1 := regexMatchSingleChar.FindStringSubmatch(input)
			inp, _ = strconv.Atoi(match1[regexMatchSingleChar.SubexpIndex("number")])
		}

		if inp >= 0 && inp <= cont-1 {
			break
		}
		cursor.ClearLinesUp(1)
		opt = fmt.Sprintf("input must be between [0, %d) ", cont)
	}
	return mapRetIndex[inp]
}

func multipleSelector[T any, S interface{}](title string, elements *[]T, maxInp int, vars S, filter func(element T, vars S) bool) []int {
	fmt.Println(title)
	cont := 0
	mapRetIndex := make(map[int]int)
	for i, v := range *elements {
		if filter(v, vars) {
			fmt.Print("\t", cont, ": ", v, "\n")
			mapRetIndex[cont] = i
			cont++
		}
	}

	var inp []int
	var opt string
	for {
		inp = []int{}
		fmt.Print("\r" + opt + "press a number to select: ")
		r := bufio.NewReader(os.Stdin)
		input, _ := r.ReadString('\n')

		if regexMatchMultipleChars.MatchString(input) {
			success := true
			match1 := regexMatchMultipleChars.FindStringSubmatch(input)
			inputs := regexSplitWs.Split(match1[regexMatchMultipleChars.SubexpIndex("args")], -1)

			if maxInp != -1 && len(inputs) > maxInp {
				opt = fmt.Sprint("too many inputs, max is ", maxInp, " ")
				continue
			}

			// check for duplicates
			for i, v := range inputs {
				for _, v2 := range inputs[i+1:] {
					if v == v2 {
						success = false
						break
					}
				}
			}

			if !success {
				opt = "input cannot contain duplicates "
				continue
			}

			for _, v := range inputs {
				inpu, err := strconv.Atoi(v)
				if err != nil || inpu < 0 || inpu > cont-1 {
					success = false
					break
				}
				inp = append(inp, inpu)
			}

			if success {
				var ret []int
				for _, v := range inp {
					ret = append(ret, mapRetIndex[v])
				}
				return ret
			}
		}
		cursor.ClearLinesUp(1)
		opt = fmt.Sprintf("input invalid, all inputs must be between [0, %d) ", cont)
	}
}

func NSpaces ( n int ) string {
	ret := ""
	for i := 0; i < n; i++ {
		ret += " "
	}
	return ret
}

func printCharacters(chs *[]Character) {

	maxLenName := 0

	for _, char := range *chs {
		if chnamelen := len(char.Name); chnamelen > maxLenName {
			maxLenName = chnamelen
		}
	}

	maxLenClass := 0

	for _, char := range *chs {
		if chclasslen := len(idToClass(char.Class)); chclasslen > maxLenClass {
			maxLenClass = chclasslen
		}
	}

	maxLenLvl := 0

	for _, char := range *chs {
		if chlvllen := len(fmt.Sprint(char.Lvl)); chlvllen > maxLenLvl {
			maxLenLvl = chlvllen
		}
	}

	for _, char := range *chs {
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
				// character has not been hit
				HpStatus = "NOHIT"
			} else if char.Hp > int(float64(char.MaxHp)*0.66) {
				// character is lightly damaged
				// at this stage mages and other fragile classes are already almost incapacitated
				HpStatus = "DAMGD"
			} else if char.Hp > int(float64(char.MaxHp)*0.33) {
				// character is wounded
				// this mostly applies to tough classes
				// at this stage median classes like the ranger are almost incapacitated
				HpStatus = "WOUND"
			} else {
				// character is at the dead door
				// the character is basically dead
				HpStatus = "DDOOR"
			}
		}

		isAlly := ' '
		if char.Friendly {
			isAlly = '⚝'
		}

		effectsString := ""

		for key, val := range char.Status {
			effectsString += fmt.Sprint( statusEffects[key].name, ": ", val, " " )
		}

		fmt.Printf(
			" %c  lvl %d%s | %s%s | %s%s | %s | %s\n",
			isAlly,
			char.Lvl, NSpaces(maxLenLvl-len(fmt.Sprint(char.Lvl))),
			char.Name, NSpaces(maxLenName-len(char.Name)),
			idToClass(char.Class), NSpaces(maxLenClass-len(idToClass(char.Class))),
			HpStatus,
			effectsString,
		)
	}
}
