.PHONY: build clean

OUT='iacm'

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${OUT}

clean:
	rm ${OUT}
