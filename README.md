# Frozen

Basic irc server made for a 48 hours rush to learn go.

### Trying it

Start the server:

```bash
git clone https://github.com/jjaniec/Frozen
cd Frozen

go build -o Frozen ./src
./Frozen
```

Initialize client connections:

```bash
nc 127.0.0.1 4242
> PASS password
> NICK nickname
> USER bob * * full name
```

A graphical client like [X-Chat Aqua](https://xchataqua.github.io/) is also supported

### Supported commands

#### Authentication

- `PASS <password>`: Set password
- `NICK <nickname>`: Set/change nickname
- `USER <user> <mode(unused)> <unused> <realname>`: Specify the username, hostname and realname of a new user

#### Channels

- `JOIN <#channel>`: Create/join a channel
- `PART <#channel>`: Leave a channel
- `PRIVMSG <(#channel/user)> :<text to be sent>`:
- `TOPIC <channel> [<topic>]`: View/set/change channel topic

- `NAMES`: Show members of each channel on the server
- `LIST`: List current channels and topics

#### Misc

- `NOTICE <(#channel/user)> <text>`: Broadcast a message to a channel/user
- `PING`: PONG
