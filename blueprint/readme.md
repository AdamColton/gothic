## Blueprint

The Blueprints package is only intended to be used by the blueprints project. A Blueprint is used to define a set of behavior in the service. In the most general case, a Blueprint will generate a stuct with methods to Marshal and Unmarshal for serialization as well as interface with entity storage. It will generate and populate the APIs and define the JavaScript for Lapiz to interpret the object.

### To-Do

blueprints need to return a package path (not just a package string) to allow nested packages

