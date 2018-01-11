package main

import "testing"

func TestSha1Sum(t *testing.T) {
  cases := []struct {
    in, want string
  }{
    {"a", "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"},
    {"b", "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98"},
  }

  for _, c := range cases {
    got := sha1sum(c.in)
    if got != c.want {
      t.Errorf("sha1sum(%q) == %q, want %q", c.in, got, c.want)
    }
  }
}

func TestGitCommit(t *testing.T) {
  git_commit()
}

func TestGitSha(t *testing.T) {
  in := "git commit"
  want := "f9a9b55eb2fad0e279ae656ef659a2d5bd17a466"
  got := git_sha(in)
  if got != want {
    t.Errorf("git_sha(%q) == %q, want %q", in, got, want)
  }
}

func TestGitMessage(t *testing.T) {
  commit := "tree 9ce1be1f32b50ebaa2ebe6513cc713312e1ce608\n" +
  "author graham jenson <grahamjenson@maori.geek.nz> 1515645586 +1300\n" +
  "committer graham jenson <grahamjenson@maori.geek.nz> 1515645586 +1300\n" +
  "\n" +
  "init\n"
  got := git_message(commit)
  want := "init"
  if(got != want) {
    t.Errorf("message wrong: %q", got)
  }
}
