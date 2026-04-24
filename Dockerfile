FROM node:24 AS jsbuild

# RUN mkdir /builddir
COPY . /builddir
WORKDIR /builddir/frontend
RUN corepack enable
RUN corepack prepare
RUN yarn install
RUN yarn build

FROM golang AS gobuild

# RUN mkdir /builddir
COPY . /builddir
COPY --from=jsbuild /builddir/frontend/dist /builddir/frontend/dist
WORKDIR /builddir
RUN go build

FROM debian
COPY --from=gobuild /builddir/hobbit-tracker /bin/hobbit-tracker
ENTRYPOINT ["/bin/hobbit-tracker"]
