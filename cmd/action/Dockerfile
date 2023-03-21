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

FROM golang@sha256:57dbdd5c8fe24e357b15a4ed045b0b1607a99a1465b5101304ea39e11547be27 AS builder
WORKDIR /src/
COPY . /src/
ARG VERSION
ARG COMMIT
ARG DATE
ENV VERSION=${VERSION}
ENV COMMIT=${COMMIT}
ENV DATE=${DATE}
RUN CGO_ENABLED=0 go build -a -trimpath -ldflags="\
    -w -s -X main.version=$VERSION \
	-w -s -X main.commit=$COMMIT \
	-w -s -X main.date=$DATE \
	-extldflags '-static'" \
    -mod vendor -o app cmd/action/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/app /app
ENTRYPOINT ["/app"]
