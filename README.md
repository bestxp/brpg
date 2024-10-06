## установка зависимостей

[ebitenengine doc](https://github.com/hajimehoshi/ebiten?tab=readme-ov-file#platforms)

## Запустить сервер

```
cd cmd/server
go run *go
```

## Запустить игру

###### Не будет работать без запущенного сервера

```
cd cmd/client
go run main.go
```

## Спрайты

[DungeonTilesetII](https://0x72.itch.io/dungeontileset-ii)

---

## Команды

```
protoc --go_out=. *.proto
sh svg2icns.sh icon.svg
```
