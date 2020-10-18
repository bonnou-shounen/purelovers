# purelovers

cli for [purelovers](https://www.purelovers.com).

## install

```bash
$ go get -u github.com/bonnou-shounen/purelovers/cmd/purelovers
```

## usage

```bash
$ export PURELOVERS_LOGIN=xxxx
$ export PURELOVERS_PASSWORD=xxxx

$ purelovers dump fav casts > casts.txt
$ vim casts.txt  # edit order
$ pruelovers restore fav casts < casts.txt

$ purelvoers dump fav shops > shops.txt
# same way
```
