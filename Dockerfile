FROM golang:1.20.4-bullseye AS builder
WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./*.go ./
RUN go build -o /out/gcloud-run-deploy-multi ./...

# ---

FROM gcr.io/distroless/cc-debian11

COPY --from=builder /out/gcloud-run-deploy-multi /bin/gcloud-run-deploy-multi

ENTRYPOINT ["/bin/gcloud-run-deploy-multi"]
