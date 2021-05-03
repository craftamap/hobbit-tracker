FROM node:15 AS jsbuild

# RUN mkdir /builddir
COPY . /builddir

RUN yarn --cwd /builddir/frontend install
RUN yarn --cwd /builddir/frontend build



FROM golang:alpine AS gobuild

# RUN mkdir /builddir
RUN apk add --update gcc musl-dev
COPY . /builddir
COPY --from=jsbuild /builddir/frontend/dist /builddir/frontend/dist
WORKDIR /builddir
RUN go build



FROM alpine
COPY --from=gobuild /builddir/hobbit-tracker /bin/hobbit-tracker
ENTRYPOINT ["/bin/hobbit-tracker"]
