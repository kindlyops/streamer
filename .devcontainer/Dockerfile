# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.162.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT="1"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

# Install Bazel
ARG BAZEL_VERSION=3.6.0
ARG BAZEL_DOWNLOAD_SHA=dev-mode
RUN curl -fSsL -o /tmp/bazel-installer.sh https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh \
  && ([ "${BAZEL_DOWNLOAD_SHA}" = "dev-mode" ] || echo "${BAZEL_DOWNLOAD_SHA} */tmp/bazel-installer.sh" | sha256sum --check - ) \
  && /bin/bash /tmp/bazel-installer.sh --base=/usr/local/bazel \
  && rm /tmp/bazel-installer.sh