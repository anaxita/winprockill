build-prod:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o winprockill.exe .