for s in foo bar baz ; do curl -d"{\"s\":\"$s\"}" localhost:8080/uppercase ; done
