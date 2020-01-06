# Kompetensdag Umeå 2020-01: Golang!

Detta repo innehåller instruktioner och källkod för kompetensdagen om Golang.

Vi kommer att börja med en introduktion av språket och därefter jobba med ett övningsprojekt. Men innan vi börjar behöver man ha utfört samtliga steg av förberedelserna nedan.

# Förberedelser

* [Installera Go][1] för ditt OS (senaste versionen)
* Bestäm var du vill lägga dina Go-projekt och sätt `GOPATH`

  Golang lägger all kod (även dina Git-repo), binärer, och andra paket i en enda katalog och den bestämmer man med `GOPATH`. Om man inte väljer något används `<hemkatalog>/go`, men om man t ex har en arbetskatalog `<hemkatalog>/Workspace` där man föredrar att lägga alla sina repon kan man ange `<hemkatalog>/Workspace/go`. Man kan om man vill lägga den i en katalog som inte heter `go` också, t ex `/foo/bar`.

* Bekräfta att Go är korrekt konfigurerat genom följande:

  * Skapa katalogen `<GOPATH>/src/hello`
  * I den katalogen, lägg in följande i `hello.go`:

  ```golang
  package main

  import "fmt"

  func main() {
      fmt.Printf("hello, world\n")
  }
  ```

  * Kompilera programmet genom att gå till katalogen du skapat och köra `go build`
  * Notera att du bör ha fått en binär `hello` i samma katalog
  * Exekvera nu programmet genom att köra `./hello`

* Konfigurera valfri editor eller IDE. T ex:

  * `VS Code` med `vscode-go`
  * `Atom` med `go-plus` (och gärna `go-debug`)
  * `Sublime` med `GoSublime` (finns också `sublime-build`)
  * `Eclipse` med `GoClipse`
  * `Vim` med `vim-go`
  * Jetbrains `GoLand` (tyvärr endast betald)
  * `Idea` med `Go` plugin (endast tillgänglig i Ultimate)

# Introduktion

Som introduktion till språket använder vi oss av [The Little Go Book][2].

Se även [Effective Go][3] som är en utmärkt referens.

# Övningsprojekt: en chat-klient

Detta projekt går ut på att vi skall bygga en klient till ett litet chat-system som har förberetts för kursen. Vi kommer att introducera den på plats, men nedan finns även lite nödvändig dokumentation.

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

[1]: https://golang.org/dl/
[2]: https://github.com/karlseguin/the-little-go-book/blob/master/en/go.md
[3]: https://golang.org/doc/effective_go.html
