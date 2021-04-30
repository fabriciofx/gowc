#!/bin/bash
for f in $(seq 1 100); do curl http://metaphorpsum.com/paragraphs/$(shuf -i 2-10 -n 1) > "dataset/file-$f.txt"; done
