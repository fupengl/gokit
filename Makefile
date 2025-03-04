.PHONY: test test-race test-coverage bench clean

# 运行所有测试
test:
	go test -v ./...

# 运行带竞态检测的测试
test-race:
	go test -race -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 运行性能测试
bench:
	go test -bench=. -benchmem ./...

# 清理生成的文件
clean:
	rm -f coverage.out coverage.html

# 运行特定包的测试
# 使用方式: make test-pkg PKG=ptr
test-pkg:
	go test -v ./$(PKG)/...

# 运行特定包的性能测试
# 使用方式: make bench-pkg PKG=ptr
bench-pkg:
	go test -bench=. -benchmem ./$(PKG)/... 