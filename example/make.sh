#!/usr/bin/env sh
# system functions
basename() {
    # Usage: basename "path" ["suffix"]
    local tmp
    tmp=${1%"${1##*[!/]}"}
    tmp=${tmp##*/}
    tmp=${tmp%"${2/"$tmp"}"}
    printf '%s\n' "${tmp:-/}"
}

lstrip() {
    # Usage: lstrip "string" "pattern"
    printf '%s\n' "${1##$2}"
}

WORK_DIR=$(pwd)
COMPANY_NAME=webdevelop-pro
SERVICE_NAME=go-logger
REPOSITORY=webdevelop-pro/go-logger

init() {
  GO_FILES=$(find . -name '*.go' | grep -v _test.go)
  PKG_LIST=$(go list ./... | grep -v /lib/)
}

build() {
  go build -ldflags "-s -w -X main.repository=${REPOSITORY} -X main.revisionID=${GIT_COMMIT} -X main.version=${BUILD_DATE}:${GIT_COMMIT} -X main.service=${SERVICE_NAME}" -o ./web/app ./web/*.go && chmod +x ./web/app
  go build -ldflags "-s -w -X main.repository=${REPOSITORY} -X main.revisionID=${GIT_COMMIT} -X main.version=${BUILD_DATE}:${GIT_COMMIT} -X main.service=${SERVICE_NAME}" -o ./cli/app ./cli/*.go && chmod +x ./cli/app
}

case $1 in

install)
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
  go install golang.org/x/lint/golint@latest
  go install github.com/lanre-ade/godoc2md@latest
  go install github.com/securego/gosec/v2/cmd/gosec@latest
  go install github.com/swaggo/swag@latest
  if [ -d ".git" -a -d ".git/hooks" ]
  then
    rm .git/hooks/pre-commit 2>/dev/null;
    ln -s etc/pre-commit .git/hooks/pre-commit
  fi
  ;;

lint)
  golangci-lint -c .golangci.yml run
  ;;

run-web)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  BUILD_DATE=$(date "+%Y%m%d")
  build && ./web/app
  ;;

run-cli)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  BUILD_DATE=$(date "+%Y%m%d")
  build && ./cli/app
  ;;

build)
  build
  ;;

help)
  cat make.sh | grep "^[a-z-]*)"
  ;;

*)
  echo "unknown $1, try help"
  ;;

esac
