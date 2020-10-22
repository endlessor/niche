run:
	go run main.go

run-seed:
	@cd seed && go run main.go

run_prod:
	export GIN_MODE=release && ./niche

build_for_linux:
	GOOS=linux GOARCH=amd64 go build -o niche
	@cd seed && go build
	@echo "Build success!"

seeder:
	@echo "Seeding get started..."
	@cd seed && ./seed

python_seeder:
	@echo "Seeding python aws..."
	@cd python && python3 py_aws.py