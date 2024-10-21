FROM chromedp/headless-shell:latest

RUN apt-get update; apt install dumb-init -y

RUN apt-get -qq install -y --no-install-recommends ca-certificates curl

ENTRYPOINT ["dumb-init", "--"]

WORKDIR /app

COPY bot /app/bot

RUN chmod +x /app/bot

ENTRYPOINT ["dumb-init", "--"]

CMD ["/app/bot","run"]

