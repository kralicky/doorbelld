FROM golang:1.16 as builder

WORKDIR /workspace
COPY . . 
RUN CGO_ENABLED=0 go build .

FROM gcr.io/distroless/static:nonroot 

WORKDIR /
COPY --from=builder /workspace/doorbelld .
USER 65532:65532

ENTRYPOINT ["/doorbelld"]