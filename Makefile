t00ls_checkin:
	CGO_ENABLED=0 go build -ldflags='-s -w' cmd/checkin/t00ls_checkin.go
clean:
	rm -f t00ls_checkin *.exe
