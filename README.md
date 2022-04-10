# secretbox

## What is secretbox?

Secretbox allows you to send notes which will self-destruct upon being read

## Why would I use it?

To keep sensitive info out of your emails and chats

## Requierements

Implemented stores for secret notes:
- in-memory db (default)
- redis

Implemented stores for sessions:
- in-memory db (default)
- redis

## Install

- download binary and start it with systemd job (TODO: embed files into binary + rewrite systemdjob)
- start as docker (TODO: write docker compose files + docker building CI task)
- deploy to kubernetes (TODO: write kube-manifests)

### Credits

Inspired by those guys:

- <a class="msg" href="https://onetimesecret.com/">One-Time Secret</a><br>
- <a class="msg" href="https://lets-go.alexedwards.net/">Alex Edwards "Let's Go!"</a>
