run:
	@cd $(root_path) && go run ./main.go serve

serve:
	@echo "fasten your belts..."
	@cd $(root_path) && ./$(app_name) serve

build:
	@cd $(root_path) && go build -o $(app_name) main.go
	@echo "project build successfully!"

start: clean build serve

clean:
	@rm -f  $(root_path)/$(app_name)
	@echo "project cleaned!"

help:
	@cd $(root_path) && go run ./main.go

app_name:= "websocket"
root_path:= ./src