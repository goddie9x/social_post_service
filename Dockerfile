FROM alpine:latest
RUN apk add --no-cache file
WORKDIR /app
COPY ./main ./.env .
EXPOSE 3005
CMD ["./main"]