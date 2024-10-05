run-web:
	go run github.com/hajimehoshi/wasmserve@latest .

run-wsl:
	GOOS=windows go run .

run:
	go run .