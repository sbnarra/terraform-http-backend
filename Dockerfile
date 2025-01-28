FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine AS go
ARG TARGETOS
ARG TARGETARCH

WORKDIR /
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o terraform-http-backend cmd/server/main.go

FROM scratch
ARG CREATED
ARG REVISION

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/terraform-http-backend"
LABEL org.opencontainers.image.title "terraform-http-backend"
LABEL org.opencontainers.image.description "Terraform HTTP Backend"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/terraform-http-backend"
LABEL org.opencontainers.image.documentation "https://github.com/sbnarra/terraform-http-backend"
LABEL org.opencontainers.image.created ${CREATED:-unset}
LABEL org.opencontainers.image.version ${REVISION:-unset}
LABEL org.opencontainers.image.revision ${REVISION:-unset}
LABEL org.opencontainers.image.base.name scratch

COPY --from=go /terraform-http-backend /terraform-http-backend
ENTRYPOINT ["/terraform-http-backend"]

EXPOSE 9944
VOLUME /data