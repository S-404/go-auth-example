This is an example implementation of an auth server to demonstrate the jwt handling in Godot project

https://github.com/S-404/godot-auth-example

# config

```cp ./configs/config.example.yaml ./configs/config.yaml```

# migrations

```go run cmd/migrate/main.go up```

```go run cmd/migrate/main.go down```


# docs

 <a href="https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format">Declarative Comments Format</a>
 

### generate


```swag init --parseDependency --parseInternal -g cmd/main.go``` 

if command 'swag' not found: `export PATH=$PATH:$(go env GOPATH)/bin` and retry

### browse

```http://localhost:8082/swagger/index.html```

# start

```go run cmd/main.go```
