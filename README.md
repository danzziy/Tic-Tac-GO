# Tic-Tac-GO

Tic-Tac-GO is a website with a GO backend and a Javascript frontend where people can play 
tic-tac-toe against there friends, family, and random people online. The vision for this site is to
have users who enter the website be presented with two choices; to play with a friend or play with
a stranger. Once a choice has been made, the user will either join an available room when playing 
with anotheruser or they will create a room to join with their friend. Once both parites are in the
room, the regular game of tic-tac-toe will begin. Once a player wins, both players will be
presented the option to either have a rematch or exit the game. If rematch is selected by both 
players, then a new game will commence within the same room. However, if either of them hits exit,
the game ends and the room is destroyed.  

## Components

The application will conist of three main components:
1) Javascript frontend for users to interact with.
2) Golang backend where user inputs will be captured and input into a database.
3) Redis database to manage game rooms.

```mermaid
graph TD;
    Javascript-->Backend;
    Backend-->Database;
```

## Redis Database

A redis database was chosen becuase I used it at work and that is what I am familiar with LMAO.
It's atomic nature will definitely be useful at the very least.

The redis database will be used to track the creation of two different types of game rooms, 
_**public**_ and _**private**_ rooms. Each room will have their own data structures associated with
it.

### Public Rooms

List of 

Associated room hash
```
PUBLIC:games:empty [<uuid1>, <uuid2>, <uuid3>]
```

```
PUBLIC:game:<uuid> full <T/F> score 
```

### Private Rooms


```
PRIVATE:game:<uuid>
```