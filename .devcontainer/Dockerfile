#trivy:ignore:DS026#trivy:ignore:DS026
FROM mcr.microsoft.com/dotnet/sdk:8.0 as azure-functions-core-tools

# Dotnet caching - see: https://github.com/dotnet/dotnet-docker/issues/3353#issuecomment-996013919
RUN --mount=type=cache,id=nuget,target=/root/.nuget/packages \ 
    git clone https://github.com/amih90/azure-functions-core-tools.git /tmp/azure-functions-core-tools
WORKDIR /tmp/azure-functions-core-tools
RUN git checkout amih90/linux-arm64 \
    && dotnet publish src/Azure.Functions.Cli/Azure.Functions.Cli.csproj --runtime linux-arm64 --output /opt/azure-functions-core-tools

# FROM golang:1.20-bookworm as final
FROM mcr.microsoft.com/devcontainers/go:1.20-bookworm as final

ARG TARGETARCH

ARG USERNAME=codespace
ARG USER_UID=1001
ARG USER_GID=$USER_UID

# APT caching - see: https://vsupalov.com/buildkit-cache-mount-dockerfile/
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
        sudo 

# Create a non-root user to use - see: https://code.visualstudio.com/remote/advancedcontainers/add-nonroot-user#_creating-a-nonroot-user
RUN groupadd --gid $USER_GID $USERNAME \
&& useradd --uid $USER_UID --gid $USER_GID -m --shell /bin/bash ${USERNAME} \
&& echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
&& chmod 0440 /etc/sudoers.d/${USERNAME}

# Install azure function core tools
# Not nativily supported for linux/arm64 see: https://github.com/Azure/azure-functions-core-tools/issues/3112
COPY --from=azure-functions-core-tools /opt/azure-functions-core-tools /opt/azure-functions-core-tools
RUN ln -s /opt/azure-functions-core-tools/func /usr/local/bin/func \
&& chmod +x /opt/azure-functions-core-tools/func \
&& rm -f azure-functions-cli.zip

USER ${USERNAME}