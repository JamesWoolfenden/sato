FROM alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

RUN apk --no-cache add build-base git curl jq bash
RUN curl -s -k https://api.github.com/repos/JamesWoolfenden/sato/releases/latest | jq '.assets[] | select(.name | contains("linux_386")) | select(.content_type | contains("gzip")) | .browser_download_url' -r | awk '{print "curl -L -k " $0 " -o ./sato.tar.gz"}' | sh
RUN tar -xf ./sato.tar.gz -C /usr/bin/ && rm ./sato.tar.gz && chmod +x /usr/bin/sato && echo 'alias sato="/usr/bin/sato"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
