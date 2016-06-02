## Gothic
For now, this is just a test to see if I can get the idea off the ground.

Gothic is a web development library. A Gothic project will have two parts, the blueprint and the service. The blueprint is a program that exists only to define the service. When the blueprint is executed, it generates much of the boilerplate code that the service needs.

### To-Do
Need a way to manage buckets (and possibly sharding) in entity.

### Validation

Simple: A single parameter can be statically validated (exists, number in range, valid email)

Resource: Parameter must be validated against a resource (object store, other parameters)

### JavaScript
* export description for Lapiz (or, whatever)
* consume description
* generate views in Lapiz
* export validation to Lapiz
* replicate serialize logic in JS