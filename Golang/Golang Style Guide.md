# Golang Style Guide

## Workspace
GOPATH, , import error? [Read this first](https://golang.org/doc/code.html)

## Naming

1. **Use camelCase** with mixedCap.
2. NO underscore.
3. Prefer short rather than long.
	- `zk` NOT `zookeeper`
4. Be consistent with abbreviations.
	- either `url` or `URL`, DO NOT use `url` here and `URL` there
5. Name should not contain type information.
	- `buffers` NOT `bufferList`
6. Make use of your context.
	- `session.New` NOT `session.NewSession`

### Package Names
1. Lower case, single-word, no need for underscores or mixedCaps.
2. DON'T worry about collisions, use whatever suits best.
	- In **rare** case of collision the importer can choose a different local name.


### Constants

1. Use iota


### Getters

1. NO `getXXX`. e.g. `Owner()` NOT `GetOwner()`.
2. Setter like `SetOwner()` is OK.


### Interfaces

1. NO `I` prefix.
2. Method name plus -er.
	- `Reader`
	- `Writer`
	- `Closer`


### Errors
1. Prefix with `err` or `Err`
	- `ErrNotFound` not `NotFound` nor `ErrorNotFound`
1. Error strings should NOT end with **punctuation**.
	- `errors.New("division by zero")` NOT `"... zero!"`


### Receiver names

1. Keep it **SHORT**, often one or two letters is sufficient.
	- It will appear in nearly every line of the method.
2. Keep it **consistent** too.
	- If you use `c` for `Client` DO NOT use `cl` in another method.

3. **DO NOT** use something like `this` `self` `me`.

_This is common in some other OOP languages, like Python, Java etc, because the line between member and non-member of the receiver(usually a class) is very clear: the members are just within the receiver, either wrapped with `{}` or indentation(Python). In go the line is blurred, thus using these generic names will reduce the readbility of your code. e.g. `self.send` and `cl.send`._