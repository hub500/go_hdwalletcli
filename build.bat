SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o hdwalletCli main.go


SET GOOS=windows
SET GOARCH=amd64
go build -o hdwalletCli.exe main.go