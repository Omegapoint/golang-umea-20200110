# golang-umea-20200110
Repo för kompetensdagen 20200110 i Umeå


## Docs

This section describes how to integrate with the system

### Name server
Holds an index of registered clients. Supported operations are currently: list clients, register client.
When a client is registered it is instantly added to the index. I.e. after a client is added to the index it will
instantly be present in the list of clients. Connected clients will have to continuously update their status.
Every five seconds, the index is pruned and clients with a connected time > 5sec will be removed from the index.
I.e. clients will have to periodically, with interval < 5sec, update their status again. A clients presence in the
list is no guarantee that it will be possible to connect to it.

---

`GET/clients`

Lists the connected endpoints. The `connected` attribute determines the last time the client announced its present.

Example response:
```
{
    "82b5268c-8231-4563-9c26-be27ea7a7abe": {
        "ip": "192.1.1.1",
        "port": 1234,
        "connected": "2019-12-19T21:20:49.49786+01:00",
        "name": "linus"
    },
    "b79534f4-810b-454a-a52e-fd4b098dca7a": {
        "ip": "192.1.1.1",
        "port": 1234,
        "connected": "2019-12-19T21:21:02.306181+01:00",
        "name": "linus"
    }
}
```

---

`POST/client`

Registers a client to the index. 

Example request body:

```
{
    "ip": "192.1.1.1",
    "port": 1234,
    "name": "linus"
}
```

Example response:

```
{
   "id": "f2ab3b6e-3284-49d4-b856-8e0939270566"
}
```

---

`PATCH/client`

Tells the name server that this client is still connected and updates the 
connected timestamp.

Example request:

```
{
    "id": "fc72368b-1f51-455b-81fc-0421fdf8666a"
}
```

Example response:

```
{
    "82b5268c-8231-4563-9c26-be27ea7a7abe": {
        "ip": "192.1.1.1",
        "port": 1234,
        "connected": "2019-12-19T21:20:49.49786+01:00",
        "name": "linus"
    },
    "b79534f4-810b-454a-a52e-fd4b098dca7a": {
        "ip": "192.1.1.1",
        "port": 1234,
        "connected": "2019-12-19T21:23:02.306181+01:00",
        "name": "linus"
    }
}
```