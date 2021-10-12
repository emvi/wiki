FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/emviwiki
COPY backend ./backend
COPY shared ./shared
COPY template /template
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# build server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/emviwiki/backend/main.go
RUN mkdir /app && \
	mv main /app/server && \
	mv backend/schema /schema && \
	rm -r /go

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
COPY --from=build /template /template
COPY --from=build /schema /schema
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    mkdir /logs && \
    mkdir /files && \
    chown -R appuser:appuser /app && \
    chown -R appuser:appuser /logs && \
    chown -R appuser:appuser /files
USER appuser

# default config
ENV EMVI_WIKI_HOST=0.0.0.0:4001
ENV EMVI_WIKI_LOGLEVEL=info
ENV EMVI_WIKI_ALLOWED_ORIGINS=*
ENV EMVI_WIKI_TLS_ENABLE=false
ENV EMVI_WIKI_BACKEND_DB_PORT=5432
ENV EMVI_WIKI_BACKEND_DB_SCHEMA=emviwiki
ENV EMVI_WIKI_BACKEND_DB_SSL_MODE=disable
ENV EMVI_WIKI_MAIL="Emvi Team <noreply@emvi.com>"
ENV EMVI_WIKI_SMTP_PORT=1025
ENV EMVI_WIKI_SMTP_SERVER=localhost
ENV EMVI_WIKI_WEBSITE_HOST=https://emvi.com
ENV EMVI_WIKI_AUTH_HOST=https://auth.emvi.com
ENV EMVI_WIKI_FRONTEND_HOST=https://emvi.com
ENV EMVI_WIKI_HTTP_WRITE_TIMEOUT=5
ENV EMVI_WIKI_HTTP_READ_TIMEOUT=5
ENV EMVI_WIKI_TEMPLATE_DIR=/template/backend/*
ENV EMVI_WIKI_MAIL_TEMPLATE_DIR=/template/mail/*
ENV EMVI_WIKI_MIGRATE_DIR=/schema
ENV EMVI_WIKI_STORE_PATH=/files

# run
EXPOSE 4001
VOLUME ["/logs", "/files"]
ENTRYPOINT ["/app/server"]
