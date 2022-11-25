package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"atomicgo.dev/cursor"
	_ "github.com/gosuri/uilive"
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
			fmt.Print("\t", cont, ": ", v, "\n" )
			mapRetIndex[cont] = i
			cont++
		}
	}

	var inp int
	var opt string
	for {
		inp = -1
		fmt.Print("\r"+opt+"press a number to select: ")
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
			fmt.Print("\t", cont, ": ", v, "\n" )
			mapRetIndex[cont] = i
			cont++
		}
	}

	var inp []int
	var opt string
	for {
		inp = []int{}
		fmt.Print("\r"+opt+"press a number to select: ")
		r := bufio.NewReader(os.Stdin)
		input, _ := r.ReadString('\n')

		if regexMatchMultipleChars.MatchString(input) {
			success := true
			match1 := regexMatchMultipleChars.FindStringSubmatch(input)
			inputs := regexSplitWs.Split( match1[regexMatchMultipleChars.SubexpIndex("args")], -1 )

			if maxInp != -1 && len(inputs) > maxInp {
				opt = fmt.Sprint("too many inputs, max is ", maxInp, " ")
				continue
			}

			// check for duplicates
			for i, v := range inputs {
				for _, v2 := range inputs[i+1:]{
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