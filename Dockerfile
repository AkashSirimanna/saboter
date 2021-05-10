FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/saboter

FROM scratch
COPY --from=build /bin/saboter /bin/saboter
ENTRYPOINT ["/bin/saboter"]
