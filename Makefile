t00ls_checkin:
	CGO_ENABLED=0 go build -ldflags='-s -w' cmd/t00ls_checkin/t00ls_checkin.go
