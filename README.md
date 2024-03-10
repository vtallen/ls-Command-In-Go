# Vincent's ls (vls)
---
This is a clone of the ls command written in go. I am using this as a way to learn the language by writing code instead of reading a bunch of tutorials

# Instalation
---
Dependencies: golang, make

run ```make install``` to build the program and move it to /usr/bin

### USAGE
---
Flags can be written as a single word with a dash in front of them or as induvidual flags. The filename must come after the flags

Examples:
    * vls <path>
    * vls -lah <path>
    * vls -l -a -h <path>
    * vls -lah
    * vls -l -a -h
  -G    Disable colorized output
  -R    List subdirectories recursively
  -S    Sort by file size
  -a    Show hidden files
  -h    Print sizes in human readable format
  -i    Print the inode number of each file
  -l    Use long listing format
  -r    Reverse the order of sort
  -t    Sort by modification time


