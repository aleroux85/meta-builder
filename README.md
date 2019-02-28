# Application Meta Builder
Build applications by copying, parsing and bundling files from a meta-source directory as described by a configuration data file.

## Meta Programming
Programming on meta file combines the benefits of configuretion files and template files to allow rapid development of applications and improved maintainability without sacrificing customisation or adding abstractions. Programming is done in a meta directory on files and templates that are used to generate the normal application source code. A configuration file contains all the data used for template parsing and file placement specification. Finally, commands can be executed on the generated files to compile, transpile or perform any other processes.

Meta programming add the benefits of:
- having a reduced code base since templates can be shared by multiple files,
- improve development speed with the ability to configure applications,
- and allows quick switching of the generated source code between different environments.

A watch mode forms a critical part of the builder. This allows the programmer to work in the meta source directory instead of the regular source directory and having changes pull through seamlessly into the entire application.
