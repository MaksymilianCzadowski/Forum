# Forum

***

This project consists in creating a web forum that allows :

* communication between users.
* associating categories to posts.
* liking and disliking posts and comments.
* filtering posts.

***
## Backend
We used golang
We used SQLite to store the data in your forum (like users, posts, comments, etc.)

## Frontend
We used framework axentix and some boostrap

## Register
Instructions for user registration:

* Must ask for email
* You need a unique email .
* You need a unique username
* You need a password
* The password is encrypted when stored

## Communication

* Only registered can create posts and comments.
* The registered users can creating a post and can associate one categories to it.
* The posts and comments are visible to all users (registered or not).

## The parts of the program

***
Here is the summary 
***

1. Starting the **server**

2. Open the server in your **web browser**

3. You can have fun creating new discussion topics with other users :)

### Explanation

***
The explanation of the parts mentioned above
***

1. It suffices to __call the program__ in the classic way using a *"go run .\server.go"*.
2. Go to you web brwoser and type in the search bar : **http://localhost:8080/login** and you will see the site appear.

### Allowed packages
* All standard go packages are allowed.
* github.com/mattn/go-sqlite3
* golang.org/x/crypto/bcrypt
* github.com/satori/go.uuid

#### Examples

***
Some examples !


**starting the server**

![Image exmple1:](https://cdn.discordapp.com/attachments/740582746979696671/853706260889206814/gorun.gif)

**Go to your server (localhost)**

![Image exmple2:](https://cdn.discordapp.com/attachments/740582746979696671/853706252243566672/local.gif)
***

