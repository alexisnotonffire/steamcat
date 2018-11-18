# Steamcat
## Picks a random game from a provided Steam library
Steamcat provides a way for those of you with ever growing Steam libraries to quickly get a game from your collection, allowing you to narrow the selection using user created categories and stats rather than the categories on Steam. This uses the `sharedconfig.vdf` file contained within the following filepath:```
`{STEAM_INSTALL}/Steam/userdata/{USERS_STEAM_ID}/7/remote/sharedconfig.vdf`
```

## Setup
First things first, grab the library: ```
go get github.com/alexisnotonffire/steamcat.git
```
There is a single dependency on a third party `.vdf` parser but everything else comes from the standard library.


## Usage
Once done, the order of operations is reasonably simple. The following code snippet is annotated to explain the order of operations to return a game from a chosen category. For simplicity, errors have been ignored.

```go
// Initiates library of games
userLib, _ := steamcat.GenerateLibraryFrom("FILE_LOCATION")

// Creates a new filtered library of games containing the given tag (currently case sensitive)
taggedLib, _ := userLib.TaggedWith("Arcade")

// Selects a random game from the library provided
game := steamcat.RandomGameFrom(taggedLib)

println(game.Name)
```

## Roadmap
The only item on the roadmap currently is to expand the library function to read from more inputes than a `.vdf` file, though this is mostly focused on the idea of receiving the `.vdf` file as part of a POST request. 
