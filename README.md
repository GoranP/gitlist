# Simple tool that uses graphql to enumerate repos in organization

Export valid Github token in environemnt variable `GITHUB_TOKEN`

Clone the code https://github.com/GoranP/gitlist.git

Compile code:
```
go get
go build
```

Run binary:
```
./gitlist -h
Usage:
  gitlist --outputcsv=<file> --orgs=<file>
  gitlist --orgs=<file> --rawjson
  gitlist -h | --help
```

Processed results will be stored in CSV file (`--outputcsv`) for Excel import. 

List of organizations are in separate text file (`--orgs`).

Eg of file:
```
org 1
org 2
```

For convenience it is possible to get raw unprocessed json in console output that is result of GraphQL query with `--rawjson` flag.

