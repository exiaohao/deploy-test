image:
	GOOS=linux GOARCH=amd64 go build -o ./docker/tester ./cmd/test
	cp -r ./test-data ./docker
	cd docker && docker build -t reg.qiniu.com/hao/deploy-test:latest .
	docker push reg.qiniu.com/hao/deploy-test:latest
	rm ./docker/tester
	rm -rf ./docker/test-data
