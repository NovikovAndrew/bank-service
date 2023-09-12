# Bank service project 

### Migration

for user migration make sure

```
$ brew install golang-migrate
```

add schema 

```
migrate create -ext sql -dir db/migration -seq init_schema
```

for migrate up use command

````
migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose up
````

### sqlc instalation 

macOS
````
brew install sqlc
````

for init yaml file user command 

````
sqlc init
````

### Install mockgen

make sure mock is set

```
ls -l ~/go/bin
```

for check mockgen user command which mockgen

in my case

```
vi ~/.zshrch
```

after insert the path  go/bin 

```
export PATH=$PATH:~/go/bin
```

and

```
source ~/.zshrc
```

### For generate mockdb use command 

```
mockgen -package mockdb -destination db/mock/store.go  bank-service/db/sqlc Store 
```

### gPRC client for testing

You can use evans for testing gRPC client

Installation:

macOS
```
brew tap ktr0731/evans
brew install evans
```

Docker image 

```
docker run --rm -v "$(pwd):/mount:ro" \
    ghcr.io/ktr0731/evans:latest \
      --path ./proto/files \
      --proto file-name.proto \
      --host example.com \
      --port 50051 \
      repl
```

```
evans -r repl
```

```
evans --tls --host example.com --port -r repl
```