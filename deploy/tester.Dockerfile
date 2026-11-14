FROM rust:1.97.0-bookworm AS web
ENV PATH=/opt/flutter/bin:$PATH \
    RUSTUP_TOOLCHAIN=1.97.0 \
    PUB_CACHE=/opt/pub-cache
RUN rustup target add wasm32-unknown-unknown \
 && cargo install wasm-bindgen-cli --version 0.2.130 --locked
RUN git clone --depth 1 --branch 3.50.0 https://github.com/flutter/flutter.git /opt/flutter \
 && flutter config --no-analytics --no-cli-animations \
 && flutter precache --web
WORKDIR /src/app
COPY app/pubspec.yaml app/pubspec.lock ./
COPY app/packages ./packages
RUN flutter pub get
COPY app/ ./
RUN tool/web_assets.sh \
 && flutter gen-l10n \
 && flutter build web --release --no-web-resources-cdn --dart-define=LODESTAR_TENANT=demo

FROM golang:1.24-alpine AS server
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
ARG VERSION=tester
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/lodestar ./cmd/lodestar \
 && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/lodestarctl ./cmd/lodestarctl

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server /out/lodestar /out/lodestarctl /
COPY --from=web /src/app/build/web /web
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/lodestar"]
