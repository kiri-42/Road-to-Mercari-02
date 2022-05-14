package main

import (
	"fmt"
	// "time"
	// "strconv"
	"os"
	"io/ioutil"
	"encoding/json"
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

	for _, word := range wordList {
		fmt.Printf("%d : %s %s\n", word.Id, word.En, word.Jp)
	}

	// var c int
	// go func() {
	// 	for {
	// 		fmt.Println("abc")
	// 		fmt.Print("-> ")
	// 		var input string
	// 		fmt.Scan(&input)
	// 		if input == "abc" {
	// 			c++
	// 		}
	// 	}
	// }()

	// select {
	// case <-time.After(5 * time.Second):
	// 	fmt.Println("\nTime's up! Score: " + strconv.Itoa(c))
	// }
}
