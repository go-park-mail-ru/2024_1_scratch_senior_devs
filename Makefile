cover:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
    grep -v -e 'mock.go' -e 'docs.go' -e '_easyjson.go' -e 'gen_notes.go' -e '.pb.go' -e 'gen.go' coverprofile_.tmp > coverprofile.tmp ; \
    rm coverprofile_.tmp ; \
    go tool cover -html coverprofile.tmp -o ../heatmap.html; \
    go tool cover -func coverprofile.tmp

test:
	go test ./...

lint:
	golangci-lint run --config=.golangci.yaml 
	
swagger:
	swag init -d cmd/main/ --parseInternal --output internal/pkg/docs --parseDependency

generate:
	go generate ./...
