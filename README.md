## Goroutines

An example go application for demonstration how to concatenate goroutines responses in same slice and print result at end.

### How works

- Search organizations members using the endpoint below:
https://api.github.com/orgs/:org/members


- In goroutines, get ID of each returned user using the endpoint below:
https://api.github.com/users/:username


- At end concatenate all results in same slice and print in terminal.

### Github Authentication

Uncomment line 32 in `main.go` and export variable `GITHUB_TOKEN_KEY` to use authentication in Github API.
```shell
export GITHUB_TOKEN_KEY=<YOUR_GITHUB_API_KEY>
```

### How to use

```shell
go run main.go
```


