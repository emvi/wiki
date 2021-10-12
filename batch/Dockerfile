FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/emviwiki
COPY batch ./batch
COPY shared ./shared
COPY auth ./auth
COPY dashboard ./dashboard
COPY backend ./backend
COPY template /template
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# build server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/emviwiki/batch/main.go
RUN mkdir /app && \
	mv main /app/server && \
	rm -r /go

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
COPY --from=build /template /template
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

# default config
ENV EMVI_WIKI_LOGLEVEL=info
ENV EMVI_WIKI_AUTH_DB_PORT=5432
ENV EMVI_WIKI_AUTH_DB_NAME=emviauth
ENV EMVI_WIKI_AUTH_DB_SSL_MODE=disable
ENV EMVI_WIKI_BACKEND_DB_PORT=5432
ENV EMVI_WIKI_BACKEND_DB_NAME=emviwiki
ENV EMVI_WIKI_BACKEND_DB_SSL_MODE=disable
ENV EMVI_WIKI_DASHBOARD_DB_PORT=5432
ENV EMVI_WIKI_DASHBOARD_DB_NAME=emvidashboard
ENV EMVI_WIKI_DASHBOARD_DB_SSL_MODE=disable
ENV EMVI_WIKI_MAIL="Emvi Team <noreply@emvi.com>"
ENV EMVI_WIKI_SMTP_PORT=1025
ENV EMVI_WIKI_SMTP_SERVER=localhost
ENV EMVI_WIKI_WEBSITE_HOST=https://emvi.com
ENV EMVI_WIKI_AUTH_HOST=https://auth.emvi.com
ENV EMVI_WIKI_FRONTEND_HOST=https://emvi.com
ENV EMVI_WIKI_MAIL_TEMPLATE_DIR=/template/mail/*

# run
VOLUME ["/logs"]
CMD ["/app/server"]
