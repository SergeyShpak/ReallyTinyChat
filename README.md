# ReallyTinyChat

ReallyTinyChat is a really tiny chat powered by the WebRTC technology.

Currently it works only with local connections.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### RTC-server

To run the RTC-server change to the __rtc-server__ directory and run the following commands:

```
go build && ./rtc-server
```

### RTC-client

1. cd to the __rtc-client__ directory:

```
cd rtc-client
```

2. Install the dependencies

```
npm install
```

3. To run the frontend in the development mode execute:

```
npm start
```

To run the frontend in the release mode execute:

TODO: add commands

### Prerequisites

For RTC-server:

1. Go version >= 1.9

For RTC-client:

1. npm (TODO: add minimal version)

## Usage

1. Start the RTC-server and the RTC-client frontend (see [Getting started](#getting-started) for details)
2. Open a tab in your browser on `localhost:3000`
3. Enter a username and a name of the chat room
4. Open another tab in your browser on `localhost:3000`
5. Enter a different username and the name of the chat room from step 3
6. Enjoy your web chat

## Running the tests

TODO

## Project To-Do list

This list is not comprehensive and not sorted in any way.

## Authors

Sergey Shpak

## License

MIT