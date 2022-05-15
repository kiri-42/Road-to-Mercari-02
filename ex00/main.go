package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
	"flag"
)

// 難易度を選択できるようにする
// 日本語訳を表示するかオプションで指定できる
// 英単語はランダムにして重複しないようにする
// エンターを押したらスタートするようにする

type word struct {
	En string  `json:"en"`
	Jp string  `json:"jp"`
}

func main() {
	useJp := flag.Bool("jp", false, "Flag to display Jp")
	flag.Parse()

	level, err := selectLevel()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	wordList, err := getWordList(level)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var correctCnt int
	ch := make(chan error, 1)
	go func() {
		var playCnt int
		size := len(wordList)
		useI := make([]int, 0)
		for {
			word := getWord(wordList, size, &useI)
			if *useJp {
				fmt.Println(word.Jp)
			}
			fmt.Println(word.En)
			fmt.Print("-> ")
			var input string
			fmt.Scan(&input)
			if input == word.En {
				correctCnt++
			}
			playCnt++
			if playCnt == size {
				ch<- errors.New("Error: No words to load")
			}
		}
	}()

	select {
	case err := <-ch:
		fmt.Println(err.Error())
		return
	case <-time.After(30 * time.Second):
		fmt.Println("\nTime's up! Score: " + strconv.Itoa(correctCnt))
	}
}

func selectLevel() (string, error) {
	fmt.Println("Please select a difficulty level(1 ~ 3)")
	fmt.Println("1: Junior High School level")
	fmt.Println("2: High School level(coming soon)")
	fmt.Println("3: University student and adult level(coming soon)")
	fmt.Print("-> ")
	var level string
	fmt.Scan(&level)
	switch level {
	case "1":
	case "2", "3":
		return "", errors.New("Coming soon!")
	default:
		return "", errors.New("Error: Please enter 1 ~ 3")
	}
	return level, nil
}

func getWordList(level string) ([]word, error) {
	bytes, err := ioutil.ReadFile("wordlist/level" + level + ".json")
	if err != nil {
		return nil, err
	}

	var wordList []word
	err = json.Unmarshal(bytes, &wordList)
	if err != nil {
		return nil, err
	}
	return wordList, nil
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
