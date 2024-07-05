feature="feat|chore|fix|hot-fix|test|revert|refactor|docs|perf|ci|build"
except="HEAD|master|main|develop|qa|staging"
code="GCA" # Change here!! GCA: Go Clean Architecture
number="[0-9]"
msg="$(cat $1)"
branchName=$(git rev-parse --abbrev-ref HEAD)

