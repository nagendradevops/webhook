# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=build

# the name of the container registry repository
REPO_NAME=<<image_registry_name/namespace>> ex:default-route-openshift-image-registry.apps.amoslondondev.amosonline.io/amos-config-manager-sla

# the name of the webhook binary file
BINARY_NAME=snow

# the name of the go command to use to build the binary
GO_CMD = go

# the version of the application
APP_VER = v0.0.1

# get old images that are left without a name from new image builds (i.e. dangling images)
DANGLING_IMGS = $(shell docker images -f dangling=true -q)

# produce a new version tag
set-version:
	sh version.sh $(APP_VER)

# build the terra binary in the current platform
build:
	$(GO_CMD) fmt
	export GOROOT=/usr/local/go; export GOPATH=$HOME/go; $(GO_CMD) build -o $(BINARY_NAME) -v

package_linux:
	export GOROOT=/usr/local/go; export GOPATH=$(HOME)/go; export CGO_ENABLED=0; export GOOS=linux; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_linux_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# build the terra docker image
snapshot-image:
	podman pull registry.access.redhat.com/ubi8/ubi-minimal
	podman build -t $(REPO_NAME)/$(BINARY_NAME)-snapshot:$(shell cat ./version) .
	podman tag $(REPO_NAME)/$(BINARY_NAME)-snapshot:$(shell cat ./version) $(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

snapshot-push:
	podman push $(REPO_NAME)/$(BINARY_NAME)-snapshot:$(shell cat ./version)
	podman push $(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

# deletes dangling
clean:
	docker rmi $(DANGLING_IMGS)
