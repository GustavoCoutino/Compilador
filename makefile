.PHONY: tidy build run erase

PROGRAM ?= program_tests/flujo_completo.patito

build:
	@echo 'Generando parser...'
	goyacc -o parser.go parser.y
	@echo 'Construyendo binario...'
	go build -o patito .

run: build
	@echo 'Parsing: $(PROGRAM)'
	./patito $(PROGRAM)

tidy:
	@echo 'Ordenando dependencias...'
	go mod tidy
	go mod verify
	@echo 'Formateando...'
	go fmt ./...

erase:
	@echo 'Destruyendo parser y binario...'
	rm -f y.output parser.go patito