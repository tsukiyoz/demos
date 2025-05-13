gen-protobuf: gen-goprotobuf gen-gogoprotobuf-fast gen-gogoprotobuf-faster gen-gogoprotobuf-slick
	
gen-goprotobuf:
	protoc -I ./proto submit.proto --go_out=./goprotobuf/submit

gen-gogoprotobuf-fast:
	protoc -I ./proto submit.proto --gogofast_out=./gogoprotobuf-fast/submit

gen-gogoprotobuf-faster:
	protoc -I ./proto submit.proto --gogofaster_out=./gogoprotobuf-faster/submit

gen-gogoprotobuf-slick:
	protoc -I ./proto submit.proto --gogoslick_out=./gogoprotobuf-slick/submit

benchmark: goprotobuf-bench gogoprotobuf-fast-bench gogoprotobuf-faster-bench gogoprotobuf-slick-bench

goprotobuf-bench:
	cd goprotobuf && go test -bench=.

gogoprotobuf-fast-bench:
	cd gogoprotobuf-fast && go test -bench=.

gogoprotobuf-faster-bench:
	cd gogoprotobuf-faster && go test -bench=.

gogoprotobuf-slick-bench:
	cd gogoprotobuf-slick && go test -bench=.