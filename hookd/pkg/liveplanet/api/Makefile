PROTOS_PATH = ./
GRPC_CPP_PLUGIN = grpc_cpp_plugin
GRPC_CPP_PLUGIN_PATH ?= `which $(GRPC_CPP_PLUGIN)`

default: gen-all

gen-all: gen-camera gen-core-go gen-metering gen-userprofile

gen-core: gen-core-go gen-core-cpp
gen-camera: gen-camera-go
gen-camera-internal: gen-camera-internal-go
gen-metering: gen-metering-go gen-metering-python
gen-userprofile: gen-userprofile-go

gen-camera-go:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		--govalidators_out=gogoimport=true:. \
		camera_cloud.proto camera_core.proto

gen-camera-python:
	python -m grpc_tools.protoc \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I./../../ \
		--python_out=./../../ \
		${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api/annotations.proto \
		&& \
	python -m grpc_tools.protoc \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I./../../ \
		--python_out=./../../ \
		${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api/http.proto \
		${GOPATH}/src/github.com/gogo/protobuf/gogoproto/gogo.proto \
		${GOPATH}/src/github.com/mwitkow/go-proto-validators/validator.proto && \
	cd ./proto && \
	python -m grpc_tools.protoc \
		-I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--python_out=. \
		--grpc_python_out=. \
		camera_cloud.proto camera_core.proto && \
	sed 's/import camera_core_pb2/from liveplanet.api.proto import camera_core_pb2/' camera_cloud_pb2.py > camera_cloud_pb2.py.new && \
	mv camera_cloud_pb2.py.new camera_cloud_pb2.py && \
	sed 's/import camera_core_pb2/from liveplanet.api.proto import camera_core_pb2/' camera_core_pb2_grpc.py > camera_core_pb2_grpc.py.new && \
	mv camera_core_pb2_grpc.py.new camera_core_pb2_grpc.py && \
	sed 's/import camera_core_pb2/from liveplanet.api.proto import camera_core_pb2/' camera_cloud_pb2_grpc.py > camera_cloud_pb2_grpc.py.new && \
	mv camera_cloud_pb2_grpc.py.new camera_cloud_pb2_grpc.py && \
	sed 's/import camera_cloud_pb2/from liveplanet.api.proto import camera_cloud_pb2/' camera_cloud_pb2_grpc.py > camera_cloud_pb2_grpc.py.new && \
	mv camera_cloud_pb2_grpc.py.new camera_cloud_pb2_grpc.py

gen-camera-internal-go:
	cd ./proto && \
	python -m grpc_tools.protoc \
		-I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		camera_cloud_internal.proto

gen-camera-internal-python:
	cd ./proto && \
	python -m grpc_tools.protoc \
		-I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--python_out=. \
		--grpc_python_out=. \
		camera_cloud_internal.proto && \
	sed 's/import camera_cloud_internal_pb2/from liveplanet.api.proto import camera_cloud_internal_pb2/' camera_cloud_internal_pb2_grpc.py > camera_cloud_internal_pb2_grpc.py.new && \
	mv camera_cloud_internal_pb2_grpc.py.new camera_cloud_internal_pb2_grpc.py

gen-core-go:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		--govalidators_out=gogoimport=true:. \
		camera_core.proto

gen-core-cpp:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		--grpc_out=. --plugin=protoc-gen-grpc=$(GRPC_CPP_PLUGIN_PATH) \
		--cpp_out=. \
		camera_core.proto

gen-metering-go:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		metering_service.proto

gen-metering-python:
	cd ./proto && \
	python -m grpc_tools.protoc \
		-I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--python_out=. \
		--grpc_python_out=. \
		metering_service.proto && \
	sed 's/import metering_service_pb2/from liveplanet.api.proto import metering_service_pb2/' metering_service_pb2_grpc.py > metering_service_pb2_grpc.py.new && \
	mv metering_service_pb2_grpc.py.new metering_service_pb2_grpc.py

gen-userprofile-go:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		userprofile_service.proto

gen-social-go:
	cd ./proto && \
	protoc -I/usr/local/include -I. \
		-I${PROTOS_PATH} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--proto_path $(PROTOS_PATH) \
		--gofast_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		upload_tasks.proto

