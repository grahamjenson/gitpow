# GitPOW

GitPOW secures the history of your git repository with proof-of-work. Use GitPOW with `gitpow 000001a` it will alter your current git commit's message to find a SHA that starts with `000001a`. If you do this for every commit in your history, if someone wanted to change history they would need to redo all your work.

# Build

`go build && go install`
