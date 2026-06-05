FROM golang:1.26.3-alpine3.23 AS dev

RUN addgroup -S appgroup && adduser -S appuser -G appgroup -h /home/appuser

WORKDIR /app

RUN go install github.com/air-verse/air@v1.65.3

COPY --chown=appuser:appgroup go.mod go.sum ./
RUN go mod download

COPY --chown=appuser:appgroup . .

USER appuser

EXPOSE 8080

CMD ["air"]

# prod
FROM scratch AS production

COPY --from=dev /etc/passwd /etc/passwd
COPY --from=dev /etc/group /etc/group

WORKDIR /

COPY --from=dev --chown=appuser:appgroup /app/main /main

USER appuser

EXPOSE 8080

ENTRYPOINT ["/main"]
