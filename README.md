# secretbox

## What is secretbox?

Secretbox allows you to send notes which will self-destruct upon being read or if they are not read within the TTL.

## Why would I use it?

To keep sensitive info out of your emails and chats
## Requirements

Implemented stores for secret notes:
- in-memory DB (INMEM)
- Redis >= 6.2.0 (REDIS)

Implemented stores for sessions:
- in-memory DB (INMEM)
- Redis >= 6.2.0 (REDIS)

## Options

Different options can be set as environment variables and overridden by flags.

| environment | Flag | default | required |
|---|---|---|---|
| LISTEN_ADDRESS | "addr" | :80 | false |
| COOKIE_SECRET | "cookie-key" | random-generated | false |
| SESSION_STORE_TYPE | "secrets-key" | --//-- | true |
| SECRETS_STORE_TYPE | "session-store-type" | --//-- | true |
| SECRETS_ENCRYPTION_KEY | "secrets-store-type" | --//-- | true |
| REDIS_ADDR | "redis-addr" | 127.0.0.1:6379 | false |

```
```

## Install

- start as docker (TODO: github ci and a guide)
- deploy to kubernetes (TODO: a guide)

## TODO
- embed files into binary
- add MySQL support
- add more tests
- expose a JSON API

### Credits

Inspired by those guys:

- <a class="msg" href="https://onetimesecret.com/">One-Time Secret</a><br>
- <a class="msg" href="https://lets-go.alexedwards.net/">Alex Edwards "Let's Go!"</a>
