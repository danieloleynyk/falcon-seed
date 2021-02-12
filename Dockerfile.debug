FROM golang:1.15.8

EXPOSE 40000

WORKDIR /
RUN go get github.com/go-delve/delve/cmd/dlv && \
    go get github.com/cosmtrek/air

HEALTHCHECK --interval=1s --timeout=10s --retries=5 CMD (ls __debug_bin && return 0) || return 1

CMD [ \
    "dlv", "debug", "./project", \
    "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--log" \
]
