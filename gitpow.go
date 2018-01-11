package main

import (
  "fmt"
  "crypto/sha1"
  "encoding/hex"
  "io"
  "os/exec"
  "time"
  "strings"
  "os"
)

func main() {
  target := "0000000"
    if(len(os.Args) > 1) {
      target = os.Args[1]
    }
    message := gitpow(target)
    fmt.Printf("new message: %v\n", message)
}

func gitpow(target string) string {
  date := time.Now().Format(time.RFC3339)
  git_init(date)
  commit := git_commit()
  old_message := git_message(commit)
  m := 0

  for(true) {
    new_message := fmt.Sprintf("%v target:%v pow:%v", old_message, target, m)
    new_commit := strings.Replace(commit, old_message, new_message, 1)
    sha := git_sha(new_commit)
    if(m % 100000 == 0) {
      fmt.Printf("nm:%q sha:%q\n", new_message, sha)
    }

    if(strings.HasPrefix(sha, target)) {
      git_update(date, new_message)
      return new_message
    }
    m += 1
  }
  return "nope"
}

func git_update(date string, new_message string) {
  date_exp := fmt.Sprintf("GIT_COMMITTER_DATE=%v", date)
  cmd := exec.Command("git", "commit", "--amend", "-m", new_message)
  cmd.Env = append(os.Environ(), date_exp)
  cmd.Run()
}

func git_init(date string) {
  date_exp := fmt.Sprintf("GIT_COMMITTER_DATE=%v", date)
  date_arg := fmt.Sprintf("--date=%v", date)

  cmd := exec.Command("git", "commit", "--amend", "--no-edit", date_arg)
  cmd.Env = append(os.Environ(), date_exp)
  cmd.Run()
}

func git_commit() string {
  cmd := exec.Command("git", "cat-file", "commit", "HEAD")
  out, _ := cmd.Output()
  return string(out)
}

func git_message(commit string) string {
  return commit[strings.Index(commit, "\n\n")+2:len(commit)-1]
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

