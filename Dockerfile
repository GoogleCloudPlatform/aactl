# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:buster AS build-env
WORKDIR /src/
COPY . /src/
ARG VERSION
ARG COMMIT
ARG DATE
ENV VERSION ${VERSION}
ENV COMMIT ${COMMIT}
ENV DATE ${DATE}
RUN CGO_ENABLED=0 go build -trimpath -ldflags="\
    -w -s -X main.version=$VERSION \
	-w -s -X main.commit=$COMMIT \
	-w -s -X main.date=$DATE \
	-extldflags '-static'" \
    -a -mod vendor -o aactl cmd/aactl/main.go

# The base image from chainguard seems to provide best balance of 
# size and vulnerabilities:
# - gogole/distroless:  13.3 MB, 9 vulnerabilities (all low)
# - chainguard/static:   5.7 MB, 0 vulnerabilities
FROM cgr.dev/chainguard/static:latest
COPY --from=build-env /src/aactl /
ENTRYPOINT ["/aactl"]
