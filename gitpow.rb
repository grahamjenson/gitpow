#!/usr/bin/ruby

require 'digest'
require 'time'

def git_sha(commit)
  Digest::SHA1.hexdigest("commit #{commit.length}\0#{commit}").to_s
end

def gitpow(target)
  date = DateTime.now.strftime("%Y-%m-%dT%H:%M:%S")
 `GIT_COMMITTER_DATE=#{date} git commit --amend --date=#{date} --no-edit`

  commit = `git cat-file commit HEAD`
  old_message = commit[commit.index("\n\n")+2..-2]

  m = 0
  while true do
    new_message = "#{old_message} target:#{target} pow:#{m}"
    sha = git_sha(commit.sub(old_message, new_message))
    puts "#{new_message} #{sha}" if m % (100000) == 0
    if sha.start_with?(target)
      puts "#{new_message} #{sha[0..7]}"
      `GIT_COMMITTER_DATE=#{date} git commit --amend -m "#{new_message}"`
      return new_message
    end
    m += 1
  end
  nil
end

gitpow(ARGV[0] || "0000000")
