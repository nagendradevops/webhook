#
#    Onix Config Manager - Alert Manager Webhook Receiver for Onix
#    Copyright (c) 2018-2020 by www.gatblau.org
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#    Unless required by applicable law or agreed to in writing, software distributed under
#    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#    either express or implied.
#    See the License for the specific language governing permissions and limitations under the License.
#
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=build

# the name of the container registry repository
REPO_NAME=<<image_registry_name/namespace>> ex:default-route-openshift-image-registry.apps.amoslondondev.amosonline.io/amos-config-manager-sla

# the name of the webhook binary file
BINARY_NAME=ses

# the name of the go command to use to build the binary
GO_CMD = go

# the version of the application
APP_VER = v0.0.5

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
