# GitPOW

GitPOW secures the history of your git repository with proof-of-work. Use GitPOW with `gitpow 000001a` it will alter your current git commit's message to find a SHA that starts with `000001a` (default `0000000`). This will protect against people rewriting history without putting all the work back in.

# Build

`go build && go install`
