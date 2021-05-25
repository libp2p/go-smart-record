# go-smart-record chat example
This sample project builds a chat application using smart-records. The app runs in the terminal and uses a UI to show messages from other peers. Making this example work straightforward:
- We start a smart-record server which will be responsible for keeping the records from clients messages.
- Clients can connect to different servers and start chatting in their chat room of choice by updating records in the server.

## Usage
To try the application you need to:

- Build the source
```
$ go build
```
- Start a server in one terminal (or the background)
```
$ ./chat-app -server
Run './chat-app -d /ip4/127.0.0.1/tcp/37615/p2p/12D3KooWSRvCLEGisZZMGRiCbzgwpQ4v9gc42tT2VHwZTYYKWfCq -room roomName -nick nickname' on another console to start a chat client.
```
- Start as much clients as you want in other terminals. You can copy the output of the server so you don't have to worry about figuring out the right address for the server.
```
$ ./chat-app -d /ip4/127.0.0.1/tcp/37615/p2p/12D3KooWSRvCLEGisZZMGRiCbzgwpQ4v9gc42tT2VHwZTYYKWfCq -room roomName -nick nickname
```
- Have fun chatting through smart-records! And remember, smart-records have a TTL, so messages may eventually disappear. 

