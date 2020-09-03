
test:
	@go test -cover ./...

dist:
	@gox \
		--osarch "!darwin/386" \
		--output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	rm -fr dist
