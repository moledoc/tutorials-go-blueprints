# Chat

This is a simple chat room application. 

**Source**: 'Go Programming Blueprints, Mat Ryer, p 15-25').

## Dependencies

* Go must be installed

## How to build and run

To build the executable/binary, run the following command in the root directory of this application:

```sh
go build
```

Then run the program with command

```sh
./chat
```

To test the chat room open two (or more) browser (tabs) at [http://localhost:8080](http://localhost:8080).

The application also allows to specify the address of the chat room.
For example:

```sh
./chat -addr=":3000"
```

open a chat room at http://localhost:3000 or

```sh
./chat -addr="192.168.0.1:3000"
```

at http://192.168.0.1:3000 (when allowed).
To find ones own ip address, then on linux you can use the command 

```sh
hostname -I
```
