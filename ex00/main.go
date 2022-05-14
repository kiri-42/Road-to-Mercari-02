package main

import (
	"fmt"
	"time"
	"strconv"
	"os"
	"io/ioutil"
	"encoding/json"
	"math/rand"
)

// 難易度を選択できるようにする
// 日本語訳を表示するかオプションで指定できる
// 英単語はランダムにして重複しないようにする
// エンターを押したらスタートするようにする

type word struct {
	Id int     `json:"id"`
	En string  `json:"en"`
	Jp string  `json:"jp"`
}

func main() {
	bytes, err := ioutil.ReadFile("wordlist/level1.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var wordList []word
	err = json.Unmarshal(bytes, &wordList)
	if err != nil {
		println(err.Error())
	}

	var correctCnt int
	go func() {
		var playCnt int
		size := len(wordList)
		useI := make([]int, 0)
		for {
			word := getWord(wordList, size, &useI)
			fmt.Println(word.En)
			fmt.Print("-> ")
			var input string
			fmt.Scan(&input)
			if input == word.En {
				correctCnt++
			}
			playCnt++
			if playCnt == size {
				println("Error")
			}
		}
	}()

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("\nTime's up! Score: " + strconv.Itoa(correctCnt))
	}
}

func getWord(wordList []word, size int, useI *[]int) word {
	var randI int
	for {
		rand.Seed(time.Now().UnixNano())
		randI = rand.Intn(size)
		if isUnique(randI, *useI) {
			*useI = append(*useI, randI)
			break
		}
	}
	return wordList[randI]
}

func isUnique(target int, useI []int) bool {
	for _, n := range useI {
		if target == n {
			return false
		}
	}
	return true
}
