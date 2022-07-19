# wormhole
An example todo application with ember and Golang for writing idiomatic/maintainable/structured Golang backend web applications.

### What problem does this solve?
People always advise using standard http library to create web applications rather than a lot of libraries out there, this repo explains how to start and structure your code. This is an example (CRUD) web application and not a library. Use this to bootstrap your web application (you may call this a tiny framework). It could also be helpful if you are looking for a way to connect emberjs with golang.

### Features
 - Model, view, controller
 - Filters (Middlewares)
 - Pass data from Filters to controllers
 - Error response to the requests
 - Validations
 - Router with namespaces and api versioning
 - static file access   

All the above can be easily tweaked and extended. As I said this is just an example.

### Installation
You will need [Golang] and [mysql] installed. Make sure that $GOPATH and $GOROOT are set correctly. 

**You will require the following libraries**   

[gorp]  
`go get github.com/coopernurse/gorp`   
[mysql driver]   
 `go get github.com/go-sql-driver/mysql`

**Download and install wormhole**   
`go get github.com/mantishK/wormhole`

**Create and configure database**   
Create a database named `wormhole`.    
Modify `wormhole/storage/mysql.go`. Assign variables `s.userName` and `s.pass` with your database username and password respectively.  
*Note: Modify `dbIp` if you are not running database on localhost.

**Install and run**    
Move to the directory where wormhole is downloaded ($GOPATH/src/github.com/mantishK/wormhole) and install it   
`go install`  
`$GOPATH/bin/wormhole`

That's it. You can test it by going to `localhost:8080` on your browser.


### Todo
 - User signup/login feature
 - Authorization
 - Authentication
 - A better readme explaining how to extend/tweak
 - A lot more better standards

### Contributions
Suggestions, pull requests, issues are all welcome.

### Credits
 - [coopernurse] for [gorp]
 - [mysql driver]
 - [emberjs] for [quickstart-code-sample]


[Golang]:https://golang.org/doc/install
[mysql]:http://www.mysql.com/downloads/
[gorp]:https://github.com/coopernurse/gorp
[mysql driver]:https://github.com/go-sql-driver/mysql
[coopernurse]:https://github.com/coopernurse
[quickstart-code-sample]:https://github.com/emberjs/quickstart-code-sample
[emberjs]:http://emberjs.com/
