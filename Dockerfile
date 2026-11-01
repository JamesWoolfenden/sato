FROM alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

RUN apk --no-cache add build-base git curl jq bash
ARG VERSION
RUN test -n "${VERSION}" || { echo "VERSION build arg is required, e.g. --build-arg VERSION=v1.2.3" >&2; exit 1; }
RUN curl -fsSL -o ./sato.tar.gz "https://github.com/JamesWoolfenden/sato/releases/download/${VERSION}/sato_${VERSION#v}_linux_386.tar.gz"
RUN tar -xf ./sato.tar.gz -C /usr/bin/ && rm ./sato.tar.gz && chmod +x /usr/bin/sato && echo 'alias sato="/usr/bin/sato"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
