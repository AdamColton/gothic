## Dev Notes

### Todo

* let user set primary on insert (override default behavoir)
* should probably expose more of this so users can add methods to helper and add templates

### conn
Should conn be passed in or should we rely on a package var?

### Meh...
I'm not sure I like that create returns a string.

MigrationFile should be handled differently. There should be a SetOnce method
and a Change method so it's not accidentally pointed at two different places.

Lines like:
h.Model.model.Model.Field(h.Primary)
are a sign I'm doing something wrong.