package main

import (
	// "fmt"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"golang.org/x/sync/errgroup"

	// "path/filepath"
	"strconv"
	"syscall"
	"os/signal"
	"time"
	// "context"
	// "net/http/httputil"
)

// http://i.imgur.com/z4d4kWk.jpg
// https://4.bp.blogspot.com/-2-Ny23XgrF0/Ws69gszw2jI/AAAAAAABLdU/unbzWD_U8foWBwPKWQdGP1vEDoQoYjgZwCLcBGAs/s1600/top_banner.jpg

const THREAD = 4

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	fileUrl := os.Args[1]

	size := getContentSize(fileUrl)

	thread := 0
	if hasAcceptRangesBytes(fileUrl) {
		thread = THREAD
	} else {
		thread = 1
	}

	ranges, err := makeRanges(thread, size)
	if err != nil {
		return
	}

	// シグナル処理
	go func() {
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-trap

		for i, _ := range ranges {
			tmpPath := getTmpPath(i)
			if f, err := os.Stat(tmpPath); !(os.IsNotExist(err) || f.IsDir()) {
				os.Remove(tmpPath)
			}
		}

		os.Exit(0)
	}()

	downloadFile(fileUrl, ranges)

	margeFile(fileUrl, ranges)
}

func margeFile(fileUrl string, ranges []string) {
	filepath := getFilepath(fileUrl)
	out, err := os.Create(filepath)
	if err != nil {
		fmt.println(os.Stderr, err.Error())
		return
	}
	defer out.Close()

	for i, _ := range ranges {
		tmpPath := getTmpPath(i)
		sub , err := os.Open(tmpPath)
		if err != nil {
			fmt.println(os.Stderr, err.Error())
			return
		}
		defer os.Remove(tmpPath)

		_, err = io.Copy(out, sub)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}
}

func downloadFile(fileUrl string, ranges []string) {
	eg := new(errgroup.Group)
	for i, s := range ranges {
		// fmt.Println(s)
		i := i
		s := s
		eg.Go(func() error {
			req, err := http.NewRequest("GET", fileUrl, nil)
			if err != nil {
				return err
			}

			req.Header.Set("Range", s)

			client := new(http.Client)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			tmpPath := getTmpPath(i)
			tmp, err := os.Create(tmpPath)
			if err != nil {
				return err
			}

			_, err = io.Copy(tmp, resp.Body)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "error :%v\n", err)
	}
}

func getTmpPath(i int) string {
	return "tmp" + strconv.Itoa(i + 1)
}

func hasAcceptRangesBytes(url string) bool {
	res, _ := http.Head(url)
	acceptRanges, _ := res.Header["Accept-Ranges"]
	for _, v := range acceptRanges {
		if v == "bytes" {
			return true
		}
	}

	return false
}

func getContentSize(url string) int {
	res, _ := http.Head(url)
	contentLength, _ := res.Header["Content-Length"]
	for _, v := range contentLength {
		n, _ := strconv.Atoi(v)
		return n
	}

	return 0
}

func makeRanges(n int, size int) ([]string, error) {
	if n < 0 {
		return nil, errors.New("thread error")
	}

	ranges := make([]string, 0)
	div := size / n
	start := 0
	end := div
	str := fmt.Sprintf("bytes=%d-%d", start, end)
	ranges = append(ranges, str)
	for end < size {
		start = end + 1
		end = start + div
		str := fmt.Sprintf("bytes=%d-%d", start, end)
		ranges = append(ranges, str)
	}

	return ranges, nil
}

func getFilepath(url string) string {
	var fileI int
	runes := []rune(url)
	if runes[len(url)-1] == '/' {
		return ""
	}

	for i, r := range runes {
		if r == '/' {
			fileI = i + 1
		}
	}

	return string(runes[fileI:])
}
