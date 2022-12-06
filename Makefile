APP_BIN = app/build/app

lint:
	golangci-lint run

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/main.go

clean:
	rm -rf ./app/build || true

git:
	git add .
	git commit -a -m "$m"
	git push -u origin main

gen:
	rm app/internal/pb/* || true
	protoc -I=app/internal/proto/ --go_out=app/internal/pb/ app/internal/proto/*.proto
	protoc --go-grpc_out=app/internal/pb/ app/internal/proto/*.proto -I=app/internal/proto/

update_contracts:
	go get -u github.com/i-b8o/read-only_contracts@$m

test:
	go test -p 1 ./...