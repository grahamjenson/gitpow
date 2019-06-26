package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	target := "0000000"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}
	message := gitpow(target)
	fmt.Printf("new message: %v\n", message)
}

func proover(id int, target string, date string, commit string, queries chan string, result chan string) {
	// recieve message to process
	// process message
	fmt.Printf("WORKER %v started\n", id)
	drain := false
	l := len(target)
	pre_old := commit[0 : strings.Index(commit, "\n\n")+2]

	for {
		new_message, more := <-queries
		if !more || drain {
			fmt.Printf("WORKER %v finished\n", id)
			return
		}

		new_commit := fmt.Sprintf("%v%v\n", pre_old, new_message)
		sha := git_sha(new_commit)

		if target == sha[0:l] {
			git_update(date, new_message)
			fmt.Printf("WORKER %v pushing to result\n", id)
			result <- new_message
			// drain the queue stop blocking on input
			drain = true
		}
	}
}

func gitpow(target string) (bool string) {
	date := os.Getenv("GIT_COMMITTER_DATE")
	if date == "" {
		date = time.Now().Format(time.RFC3339)
	}

	git_init(date)
	commit := git_commit()
	old_message := git_message(commit)
	workers := runtime.NumCPU() - 1

	if workers <= 0 {
		workers = 1
	}

	// Channel for numbers to be processed buffer of 5 per worker
	queries := make(chan string, workers*5)

	// results for processed numbers (requires size of workers in case of both finding an answer)
	result := make(chan string, workers)

	// Nu,ber of workers

	for w := 1; w <= workers; w++ {
		go proover(w, target, date, commit, queries, result)
	}

	m := 0
	for {
		select {
		case res := <-result:
			// closing the queue will set more to false
			close(queries)
			close(result)
			return res
		default:
			new_message := fmt.Sprintf("%v target:%v pow:%v", old_message, target, m)
			if m%1000000 == 0 {
				fmt.Printf("nm:%q\n", new_message)
			}
			queries <- new_message
			m += 1
		}
	}
	fmt.Printf("RETURNING NOTHING\n")
	return "nope"
}

func git_update(date string, new_message string) {
	date_exp := fmt.Sprintf("GIT_COMMITTER_DATE=%v", date)
	cmd := exec.Command("git", "commit", "--amend", "--allow-empty", "--no-gpg-sign", "-m", new_message)
	cmd.Env = append(os.Environ(), date_exp)
	cmd.Run()
}

func git_init(date string) {
	date_exp := fmt.Sprintf("GIT_COMMITTER_DATE=%v", date)
	date_arg := fmt.Sprintf("--date=%v", date)

	cmd := exec.Command("git", "commit", "--amend", "--no-gpg-sign", "--no-edit", date_arg)
	cmd.Env = append(os.Environ(), date_exp)
	cmd.Run()
}

func git_commit() string {
	cmd := exec.Command("git", "cat-file", "commit", "HEAD")
	out, _ := cmd.Output()
	return string(out)
}

func git_message(commit string) string {
	return commit[strings.Index(commit, "\n\n")+2 : len(commit)-1]
}

func git_sha(commit string) string {
	return sha1sum(fmt.Sprintf("commit %v\u0000%v", len(commit), commit))
}

func sha1sum(s string) string {
	hasher := sha1.New()
	io.WriteString(hasher, s)
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}
