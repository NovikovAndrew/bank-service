# Bank service project 

### Migration

for user migration make sure

$ brew install golang-migrate

migrate create -ext sql -dir db/migration -seq init_schema

for migrate up use command
migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose up

sqlc instalation 

macOS
brew install sqlc

Ubuntu
sudo snap install sqlc

Docker
docker pull kjconroy/sqlc

for init yaml file user command 
sqlc init

install mockgen

make sure mock is set

ls -l ~/go/bin

for check mockgen user command which mockgen

in my case vi ~/.zshrch

after insert the path  go/bin 

export PATH=$PATH:~/go/bin

and 

source ~/.zshrc

