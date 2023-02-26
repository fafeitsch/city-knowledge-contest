FROM golang:1.18-alpine AS BUILD_BACKEND
WORKDIR /app
COPY ./backend .
RUN go mod download
RUN go build -o contest-server ./cmd/contest-server.go

FROM node:18-alpine AS BUILD_FRONTEND
WORKDIR app
COPY ./frontend .
RUN npm install
RUN npm run build

FROM alpine:latest
COPY --from=BUILD_BACKEND /app/contest-server contest-server
RUN mkdir frontend
COPY --from=BUILD_FRONTEND /app/dist/ ./frontend/
RUN chmod a+x frontend
EXPOSE 23123
ENTRYPOINT ["/contest-server"]
CMD []
