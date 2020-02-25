OUTPUT_DIR = ./builds

tools:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr

build:
	GOARM=7 go run github.com/mitchellh/gox -os="darwin linux" -arch="386 amd64 arm" -osarch="!darwin/arm" -output "${OUTPUT_DIR}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/pkg/; for osarch in *; do (cd $$osarch; tar zcvf ../../dist/fifo_broadcaster_$$osarch.tar.gz ./*); done;

clean:
	rm -rf ${OUTPUT_DIR}
