FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/emviwiki
COPY auth ./auth
COPY shared ./shared
COPY template /template
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# build server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/emviwiki/auth/main.go
RUN mkdir /app && \
	mv main /app/server && \
	mv auth/static /app/static && \
	mv auth/schema /schema && \
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
    chown -R appuser:appuser /app && \
    chown -R appuser:appuser /secrets
USER appuser

# default config
ENV EMVI_WIKI_HOST=0.0.0.0:4000
ENV EMVI_WIKI_LOGLEVEL=info
ENV EMVI_WIKI_ALLOWED_ORIGINS=*
ENV EMVI_WIKI_TLS_ENABLE=false
ENV EMVI_WIKI_AUTH_DB_PORT=5432
ENV EMVI_WIKI_AUTH_DB_SCHEMA=emviauth
ENV EMVI_WIKI_AUTH_DB_SSL_MODE=disable
ENV EMVI_WIKI_TOKEN_PUBLIC_KEY=/secrets/token.public
ENV EMVI_WIKI_TOKEN_PRIVATE_KEY=/secrets/token.private
ENV EMVI_WIKI_MAIL=noreply@emvi.com
ENV EMVI_WIKI_SMTP_PORT=1025
ENV EMVI_WIKI_SMTP_SERVER=localhost
ENV EMVI_WIKI_AUTH_HOST=https://auth.emvi.com
ENV EMVI_WIKI_TEMPLATE_DIR=/template/auth/*
ENV EMVI_WIKI_MAIL_TEMPLATE_DIR=/template/mail/*
ENV EMVI_WIKI_MIGRATE_DIR=/schema

# run
EXPOSE 4000
VOLUME ["/secrets", "/logs"]
ENTRYPOINT ["/app/server"]
