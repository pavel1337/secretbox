# secretbox

## What is secretbox?

Secretbox allows you to send notes which will self-destruct upon being read or if they are not read within the TTL. [Example](https://secretbox.ipvl.de)

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
| SECRETS_ENCRYPTION_KEY | "secrets-key" | --//-- | true |
| SESSION_STORE_TYPE | "session-store-type" | --//-- | true |
| SECRETS_STORE_TYPE | "secrets-store-type" | --//-- | true |
| REDIS_ADDR | "redis-addr" | 127.0.0.1:6379 | false |

## Install

- start as docker with in-memory DB
```shell
docker pull pv1337/secretbox
docker run -p 8080:8080 --name secretbox pv1337/secretbox --addr :8080 \
--secrets-key your_32char_supervery_secret_key \
--session-store-type INMEM --secrets-store-type INMEM
```

- deploy to kubernetes with in-memory DB
```shell
# clone this repo or download kubernetes/ folder

kubectl create namespace secretbox
kubectl create -n secretbox secret generic \
    --from-literal=SECRETS_ENCRYPTION_KEY=yahdeoraquiekeezieJ4thi6ahn0ze9u \
    --from-literal=COOKIE_SECRET=Xui4Phoogohquaecie5ier4aaghoh1su \
    secretbox-secrets
kubectl create -n secretbox configmap  \
    --from-literal="LISTEN_ADDRESS=:8080" \
    --from-literal=SESSION_STORE_TYPE=INMEM \
    --from-literal=SECRETS_STORE_TYPE=INMEM \
    secretbox-configmap

kubectl apply -n secretbox -f kubernetes/deployment.yml
kubectl apply -n secretbox -f kubernetes/service.yml

# ingress setup is outside of the scope of this guide
# but you can port-forward and check if it is running
kubectl port-forward -n secretbox svc/secretbox-service 8080
```
## TODO
- embed files into binary
- add MySQL support
- add more tests
- expose a JSON API

### Credits

Inspired by those guys:

- <a class="msg" href="https://onetimesecret.com/">One-Time Secret</a><br>
- <a class="msg" href="https://lets-go.alexedwards.net/">Alex Edwards "Let's Go!"</a>
