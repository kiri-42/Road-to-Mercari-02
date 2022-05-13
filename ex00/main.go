package main

import (
	"fmt"
	"time"
	"strconv"
)

// 難易度を選択できるようにする
// 日本語訳を表示するかオプションで指定できる
// 英単語はランダムにして重複しないようにする
// エンターを押したらスタートするようにする

func main() {
	var c int
	go func() {
		for {
			fmt.Println("aaaa")
			fmt.Print("-> ")
			var input string
			fmt.Scan(&input)
			if input == "aaaa" {
				c++
			}
		}
	}()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("\nTime's up! Score: " + strconv.Itoa(c))
	}
}
