# Makefile для запуска всех микросервисов из директории APIGateway

# Список микросервисов
SERVICES := censors comments news gateway

.PHONY: all $(SERVICES)

all: $(SERVICES)

$(SERVICES):
	@echo "Launching the service $@"
	@(cd $@ && start cmd /c "go run cmd/main.go")

