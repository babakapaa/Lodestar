FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/lodestar ./cmd/lodestar \
 && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/lodestarctl ./cmd/lodestarctl \
 && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/lodestar-demo-idp ./cmd/lodestar-demo-idp

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/lodestar /out/lodestarctl /out/lodestar-demo-idp /
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/lodestar"]
