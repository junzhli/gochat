# Chat system
Server/Client architecture chat system implemented in Go  
The design is derived from [chatserver](https://github.com/nqbao/learn-go/tree/chat/0.0.1/chatserver), and I add more 
features based upon this

* Single chat room
* System-wide notifications (e.g. user is online/offline)
* User must set nickname to proceed
* Broadcast message to everyone
* Error handling improved

### Build

```shell script
>> make build
```

Compiled binaries output to `dist` folder

* Server
  * `gochatd`
* Client
  * `gochat`

### Install

Build and place in `$GOPATH/bin`

```shell script
>> make install
```