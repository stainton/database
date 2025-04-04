
# 可以使用“make startup”快速启动服务
.PHYON: startup
startup:
	go run main.go --db-server=localhost --db-port=3306 --db-user=root --db-passwd=961110 --db-name=testdb --addr=localhost --port=8090