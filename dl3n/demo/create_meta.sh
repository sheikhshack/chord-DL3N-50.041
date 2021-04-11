#!/bin/sh

# creates a demo file and a metadata file for it.
figlet "this is a demo file" > shared_file.txt
./dl3n writemeta shared_file.txt
