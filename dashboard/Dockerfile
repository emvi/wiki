FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/emviwiki
COPY dashboard ./dashboard
COPY auth ./auth
COPY backend ./backend
COPY shared ./shared
COPY template /template
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# build server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/emviwiki/dashboard/main.go
RUN mkdir /app && \
	mv main /app/server && \
	mv dashboard/static /app/static && \
	mv dashboard/schema /schema && \
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
    mkdir /secrets && \
    mkdir /logs && \
    mkdir /files && \
    chown -R appuser:appuser /app && \
    chown -R appuser:appuser /secrets && \
    chown -R appuser:appuser /logs && \
    chown -R appuser:appuser /files
USER appuser

# default config
ENV EMVI_WIKI_HOST=0.0.0.0:4005
ENV EMVI_WIKI_LOGLEVEL=info
ENV EMVI_WIKI_ALLOWED_ORIGINS=*
ENV EMVI_WIKI_TLS_ENABLE=false
ENV EMVI_WIKI_AUTH_DB_PORT=5432
ENV EMVI_WIKI_AUTH_DB_NAME=emviauth
ENV EMVI_WIKI_AUTH_DB_SSL_MODE=disable
ENV EMVI_WIKI_BACKEND_DB_PORT=5432
ENV EMVI_WIKI_BACKEND_DB_NAME=emviwiki
ENV EMVI_WIKI_BACKEND_DB_SSL_MODE=disable
ENV EMVI_WIKI_DASHBOARD_DB_PORT=5432
ENV EMVI_WIKI_DASHBOARD_DB_NAME=emvidashboard
ENV EMVI_WIKI_DASHBOARD_DB_SSL_MODE=disable
ENV EMVI_WIKI_TOKEN_PUBLIC_KEY=/secrets/token.public
ENV EMVI_WIKI_TOKEN_PRIVATE_KEY=/secrets/token.private
ENV EMVI_WIKI_TEMPLATE_DIR=/template/dashboard/*
ENV EMVI_WIKI_MAIL_TEMPLATE_DIR=/template/mail/*
ENV EMVI_WIKI_MIGRATE_DIR=/schema
ENV EMVI_WIKI_STORE_PATH=/files

# run
EXPOSE 4005
VOLUME ["/secrets", "/logs", "/files"]
ENTRYPOINT ["/app/server"]
