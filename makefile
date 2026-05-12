tidy:
	@echo 'Ordenando dependencias de modulos...'
	go mod tidy
	@echo 'Verificando y proveyendo dependencias de modulos...'
	go mod verify
	go mod vendor
	@echo 'Formateando archivos de .go ...'
	go fmt ./...

PROGRAM ?= tests/programa1.patito

build:
	@echo 'Generando parser...'
	goyacc -o parser.go parser.y
	@echo 'Construyendo binario...'
	go build -o patito
	@echo 'Parsing archivo de entrada: $(PROGRAM)...'
	./patito $(PROGRAM)

erase:
	@echo 'Destruyendo parser y output...'
	rm y.output parser.go patito