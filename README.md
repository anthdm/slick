# Slick
Build web applications faster with all batteries included

> Work in progress

# I really want to try this out even though its not finished yet.
Spoiler: you will encounter rough edges.

For now you will need to install templ and air manually (working on making this work with `slick install`

- [https://github.com/a-h/templ/](https://github.com/a-h/templ/)
- [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air)

Install the Slick cli
```
go install "github.com/anthdm/slick/slick@latest"
```

Create new slick project
```
slick new myapp
```

Install the project
```
cd myapp && slick install
```

Start the project
```
slick run
```

Run application in watch mode using [air](https://github.com/cosmtrek/air)
```
air
```
