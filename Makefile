
test:
	@go test -cover ./...

dist:
	@gox -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	rm -fr dist
