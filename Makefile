SHELL:=/bin/bash

release:
	@ chmod +x ./ci/release.sh
	@ ./ci/release.sh ${PWD}/version.go

release-binaries:
	@ chmod +x ./ci/add-release-assets.sh
	@ ./ci/add-release-assets.sh