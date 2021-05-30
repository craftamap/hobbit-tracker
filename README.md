# Hobbit Tracker

Hobbit Tracker is a small habit tracker built with Golang and Vue.js.

## Building the project

### Prerequisites

The backend is built in go and requires go 1.16 or later.

The frontend is built using yarn.

### Building

Make sure to built the frontend first, as the backend build process includes the
frontend into it's binary.

#### Frontend 

```
cd frontend
yarn install
yarn build
```

#### Backend

```
go get -v
go build
```

### Running

```
./hobbit-tracker -port 8080
# a sqlite file will be created called hobbits.sqlite
```

### Developing 

Currently, we use the following important libraries in the backend:

- [gorilla/mux](https://github.com/gorilla/mux)
- [gorm.io/gorm](https://gorm.io/gorm)
- [sirupsen/logrus](https://github.com/sirupsen/logrus)

To launch the backend for development, you can also use the `-disk-mode`-flag, 
which will read the files from your disk, and not from the files included in the binary.

The frontend uses Vue.js, Vue Router and Vuex, as well as some other smaller libraries.

### Contributing

Although this is a private project of mine, feel free to contribute to it.
