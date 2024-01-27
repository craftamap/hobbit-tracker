FROM node:18 AS jsbuild

# RUN mkdir /builddir
COPY . /builddir

RUN yarn --cwd /builddir/frontend install
RUN yarn --cwd /builddir/frontend build



FROM golang AS gobuild

# RUN mkdir /builddir
COPY . /builddir
COPY --from=jsbuild /builddir/frontend/dist /builddir/frontend/dist
WORKDIR /builddir
RUN go build



FROM debian
COPY --from=gobuild /builddir/hobbit-tracker /bin/hobbit-tracker
ENTRYPOINT ["/bin/hobbit-tracker"]
