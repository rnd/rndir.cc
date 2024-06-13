BUILD_DIR := build

.PHONY: dev
dev: templ server run

.PHONY: gen
gen: dep templ resume

.PHONY: server
server:
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/main \
		-ldflags="-w -s -X 'main.version=local'" \
		./cmd/site

.PHONY: run
run:
	[ -f ./env/config ] || { cp ./env/sample.config ./env/config; }
	$(BUILD_DIR)/main

.PHONY: test
test:
	CGO_ENABLED=1 \
		go test -vet=off -count=1 -race -timeout=3s ./...

.PHONY: templ
templ:
	templ fmt .
	templ generate

.PHONY: dep
dep:
	gomod2nix generate

.PHONY: resume
resume:
	xelatex -output-directory=$(BUILD_DIR) static/resume.tex
	cp $(BUILD_DIR)/resume.pdf static/resume.pdf

.PHONY: image
image:
	nix build --file ./nix/do.nix
