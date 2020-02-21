FROM golang:1.13 as builder

WORKDIR /app
COPY . .
RUN make

FROM scratch
ARG BUILD_DATE
ARG VCS_REF

COPY --from=builder /app/dist/chatbot-*-linux ./chatbot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
EXPOSE 4000
ENTRYPOINT ["./chatbot"]
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="chatbot" \
      org.opencontainers.image.authors="Anthony corbacho" \
      org.opencontainers.image.source="https://github.com/anthonycorbacho/chatbot/build" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.vendor="Anthony corbacho"