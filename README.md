# Application Meta Builder (Library)
Build applications by copying, parsing and bundling files from a meta-source directory as described by a configuration data file.

## Meta Programming
Programming on meta files combines the benefits of configuration files and template files to allow rapid development of applications and improved maintainability without sacrificing customisation or adding abstractions. Programming is done in a meta directory on files and templates that are used to generate the regular application source code. A configuration file contains all the data used for template parsing and file placement specification. Finally, commands can be executed on the generated files to compile, transpile or perform any other processes.

Meta programming add the benefits of:
- having a reduced code base since templates can be shared by multiple files,
- improve development speed with the ability to configure applications,
- and allows quick switching of the generated source code between different environments.

A watch mode forms a critical part of the builder. This allows the programmer to work in the meta source directory instead of the regular source directory and having changes pull through seamlessly into the entire application.

## Features

The builder uses the `meta.json` configuration file
and the template files (in the meta folder) to
build the project.

The builder can be executed with:

```bash
bin/meta-builder
```

to build the project. Files that does not exist yet will be created. If files exist, it will not be replaced.
The flag `-f` forces the replacement of all files and can be used as such:

```bash
bin/meta-builder -f
```

A `-w` flag can be added to put the builder into
**watch mode** where files will be automatically
updated when the meta code is changed.

The builder will transfer files from sources to destinations by:
- stepping through directories recursively
- parsing any templates available in the source specified by the "from" key.
- stepping through files
- use file key as source name if "source" field is empty
- use file key as destination name
- create destination file

##### File Data Structures



## meta.json Configuration Reference

The config file defines the project at the top level. The project structure is:

```json
{
  "name": "project-name",
  "directories": {},
}
```

File structures (FSs) are specified with the `directories` key.
The `directories` key contain key-value pairs of multiple FSs
that can each be understood as a directory in the project.

```json
{
  "directories": {
    ...
    "dir-one": {},
    "dir-two": {},
    ...
  }
}
```

A FS can contain files as key-value pairs under the `files` key.
These are the files that will be built.

```json
{
  "directories": {
    ...
    "dir-two": {
      "files": {
        ...
        "file-aaa.ext": {},
        "file-bbb.ext": {},
        ...
      }
    },
    ...
  }
}
```

By default, a file in the meta directory with the path name `file/path/file-name`
will be parsed and written to a file (named with the file key)
and placed under a directory (named with the FS key)
in the project root directory (`project-root/FS-name/file-name`).

```json
{
  "directories": {
    "one": {
      "files": {
        "aaa.ext": {},
        "bbb.ext": {}
      }
    },
    "two": {
      "files": {
        "ccc.ext": {}
      }
    },
  }
}
```
will build to:

```
./meta/one/aaa.ext -> ./one/aaa.ext
./meta/one/bbb.ext -> ./one/bbb.ext
./meta/two/ccc.ext -> ./two/ccc.ext
```


More FS option keys are:
- `from: "source-directory"` for specifying a sub-directory in the meta directory where the source file will be found.
- `dest: "destination-directory"` for modifying the destination path output files. Various options are available (consider FS key `b` in FS key `a`):
  - specifying an additional directory name or path (`dest: "c"` -> `a/c/file.ext`, `dest: "c/d"` -> `a/c/d/file.ext`)
  - using a `./` to ignore the FS key, not making a sub-directory (`dest: "./"` -> `a/file.ext`, `dest: "./c"` -> `a/c/file.ext`)
  - using a `/` to go back to the project root (`dest: "/"` -> `file.ext`, `dest: "/c"` -> `c/file.ext`)
- `copyfiles: true` to only copy file and skip parsing.
- `directories: { FS-name: {} }` for specifying child FSs. Note that child FSs act as sub-directories in the output path.

## Using as stand-alone


## Implementing your own builder

