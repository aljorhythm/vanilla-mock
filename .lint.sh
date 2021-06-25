#!/bin/sh

formatted_files=$(go fmt ./...)
echo $formatted_files

[ -n "$formatted_files" ] && echo "warning: formatted go files, please retry" && exit 1

exit 0