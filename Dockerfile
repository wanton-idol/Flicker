# FROM golang:1.19


# RUN mkdir /app
# ADD . /app
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# # Build
# RUN go build -o main .


# EXPOSE 8080

# # Run
# CMD [ "/app/main" ]


FROM golang:1.19.5@sha256:a0b51fe882f269828b63e7f69e6925f85afc548cf7cf967ecbfbcce6afe6f235 as base

# go-dev-container stage is used as a vs-code dev container (see .devcontainer/devcontainer.json for more info).
FROM base as go-dev-container

RUN apt-get update -y && \
  apt-get install -y vim zsh less && \
  apt-get install -y yamllint && \
  chsh -s $(which zsh) && \
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
  go install github.com/vektra/mockery/v2@v2.9.4 && \
  go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
  go install github.com/ramya-rao-a/go-outline@latest && \
  go install github.com/cweill/gotests/gotests@latest && \
  go install github.com/fatih/gomodifytags@latest && \
  go install github.com/josharian/impl@latest && \
  go install github.com/haya14busa/goplay/cmd/goplay@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  ( \
  GOBIN=/tmp/ go install github.com/go-delve/delve/cmd/dlv@latest && \
  mv /tmp/dlv $GOPATH/bin/dlv-dap \
  ) && \
  go install honnef.co/go/tools/cmd/staticcheck@v0.2.2 && \
  go install golang.org/x/tools/gopls@latest


# go-builder is a different stage than go-dev-container as it builds the go app.
# If they were the same stage it would be impossible to build go-dev-container if there were build errors
# FROM go-dev-container as go-builder

# WORKDIR /build

# COPY go* ./

# RUN go mod download && \
#   apt-get install -y yamllint && \
#   mkdir bin

# COPY . .

# RUN go build -o ./bin ./...

# FROM base as production


# WORKDIR /app

# RUN apt-get update --assume-yes && \
#   apt-get install --assume-yes \
#   unzip \
#   libc6 \
#   groff \
#   less && \
#   curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.30.zip" -o "awscliv2.zip" && \
#   unzip awscliv2.zip && \
#   ./aws/install && \
#   rm -f awscliv2.zip

# ARG DD_AGENT_VERSION="1:7.43.*"

# ENV DD_INSTALL_ONLY=true \
#   USE_DOGSTATSD=yes

# Some of these packages may be redundant after removing Ruby from datastore-monitor
# RUN apt-get update --assume-yes && \
#   apt-get install --assume-yes apt-transport-https && \
#   sh -c "echo 'deb https://apt.datadoghq.com/ stable 7' > /etc/apt/sources.list.d/datadog.list" && \
#   apt-key adv --recv-keys --keyserver hkp://keyserver.ubuntu.com:80 D75CEA17048B9ACBF186794B32637D44F14F620E && \
#   apt-get update --assume-yes && \
#   apt-get install --assume-yes \
#   build-essential \
#   tzdata \
#   traceroute \
#   dirmngr \
#   jq \
#   datadog-agent=$DD_AGENT_VERSION && \
#   find -name /etc/datadog-agent/conf.d/conf.yaml.default -execdir mv {} conf.yaml.disabled \;

# COPY . /app

# COPY --from=go-builder /build/bin/Supermatch bin/

# CMD ["bin/supermatch/main"]