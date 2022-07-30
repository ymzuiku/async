v=v0.0.9
git tag $v
git push --tags
go install github.com/ymzuiku/async@$v
echo "done."