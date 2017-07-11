## Dev Notes

### Todo

#### Relational Models
I need a way to setup relational models.

#### Stack Definition
Within an applicaiton, most models will need the same stack (storage, server,
net-comm and client). It would be helpful to be able to define a stack and give
it a model and get all the pieces at once.

### Broad, sweeping plan
This may be a bad idea...

#### Data Models
Like Go structs. Defines "objects" as ordered, key/value pairs. The key is the
field name and the value is the type. A model has a name, package and a version
number. The version number is 2 parts.

Objects need to support relational models. And indexing models.

There also needs to be something like inheiritance. It could be Go style
embedding. It could be classical inheiritance or JS prototypes. Not sure.

In this context is also migrations. A migration is a set of changes. A single
migration can effect several models. Migrations change the model number. Changes
that are backwards compatible (like adding a field) increment minor version,
while changes that are not increment the major version.

#### Service Models
This are the most fuzzy in my mind right now. This is fundamentally a way to
describe nodes and channels that form a network.

It should be easy to "import" something like an http-rest definition and build
on that, but the core concept should not be that.

A computational node describes something can fulfill endpoints over channels.

A channel describes a communication protocol. It can be as simple as being
named ("http", "udp", "rpc") or it can describe requirements (bi-directional,
synchronous, has headers).

An endpoint is like an API. It describes how messages will be routed within
the service.

And a message describes what the data should look like. Headers and body.

#### Compute Models
Has a package, shared with the Data Model. Within a package are functions,
methods and interfaces. Kind of fuzzy still.

Just thinking; a function/method can either describe a "contract" or prescribe
an algorithm. An interface is just like a Go interface.

#### Bringing it together
All of these will need a DSL and a gothic library. Do the lib first.

It should be possible to write this and then spit out all much of the necessary
code.

An example case; describe some models including indexing. Describe a client and
server node and communications model. Describe some validation operations on the
models. From this, you can get your database migrations, baked-in orm style
operations, models in Go and JS along with a lot of boilerplate code, including
the validation code that was defined and, finally, client and server API
libraries.

But the code can also be library dependant. So migration from one front-end
framework to another should be easier. Or migrating from sql to no-sql is
simple.

And as the application evolves, the history of the models is kept.