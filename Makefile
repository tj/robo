
test:
	@go test -cover ./...

.PHONY: dist
dist:
	@gox \
		--osarch "darwin/amd64 darwin/arm64 linux/amd64 linux/arm" \
		--output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
	rm dist/*.gz
	ls dist/* | while read i; do gzip $$i; done

clean:
	rm -fr dist
