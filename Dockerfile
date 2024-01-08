FROM node:18.19.0 AS frontend
WORKDIR /frontend
COPY ./frontend .
RUN npm ci
RUN npm run build

FROM golang:1.21.5 AS backend
WORKDIR /go/src/avia-calendar
COPY . .
COPY --from=frontend --chown=$USER:$USER /frontend/dist ./frontend/dist
RUN go build -o avia-calendar ./cmd/avia-calendar-rest

FROM ubuntu:22.04
WORKDIR /
COPY --from=backend --chown=$USER:$USER /go/src/avia-calendar/avia-calendar ./avia-calendar

EXPOSE 8796

ENTRYPOINT ["/avia-calendar"]