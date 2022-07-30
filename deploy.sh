v=v0.0.5
git tag $v
git push --tags
go install github.com/ymzuiku/async@$v
echo "done."