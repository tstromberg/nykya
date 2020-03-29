# nykya

Static blogging platform.


## Setting a root directory

If you get tired of specifying `nykya --root` to point to the site directory, you can set it as a persistent root directory:

* UNIX: `export NYKYA_ROOT=/site`
* Windows: `[System.Environment]::SetEnvironmentVariable('NYKYA_ROOT', 'C:\Users\spam1\site', [System.EnvironmentVariableTarget]::User)`


## Reference

nykya <verb> <object> <content>

## Supported examples

* `ny add thought "Where am I going with my life?` - record a new thought
* `ny add post` - open a text editor to record a new post
* `ny add post /path/to/post.md`:
   - If within the site root, append frontmatter
   - If outside of site root, copy it and append frontmatter
* `ny add image /path/to/image`
   - Same semantics as a new post

## Verbs

* add - add something
* rm - remove something
* dev - startup development webserver
* sync - resync post content
* render - write static output

## Objects

* photo - local or remote URL to JPG
* post - local or remote URL to HTML
   - FUTURE: Google Docs integration
* thought - quick inline text post

## Flags

* --root - site directory containing nykya.yaml file (config)
* --description -  