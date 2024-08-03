#!/bin/sh

if [ $# -lt 3 ]; then
    echo "Usage: $0 library old_prefix new_prefix action"
    exit 1
fi

LIBRARY=$1
WRONG_PREFIX=$2
RIGHT_PREFIX="@executable_path/../$3"
ACTION=$4

chmod u+w $LIBRARY

if [ "x$ACTION" = "xchange" ]; then
    libs="`otool -L $LIBRARY 2>/dev/null | fgrep compatibility | cut -d\( -f1 | grep $WRONG_PREFIX | sort | uniq`"
    for lib in $libs; do
        if ! echo $lib | grep --silent "@executable_path" ; then
            if echo $lib | grep --silent "${LIB_HOME}/Cellar/"; then
                fixed=`echo $lib | sed -e "s|${LIB_HOME}/Cellar/\([^/]*\)/[^/]*/|@executable_path/../Resources/opt/\1/|"`
            else
                fixed=`echo $lib | sed -e s,\${WRONG_PREFIX},\${RIGHT_PREFIX},`
            fi
            echo  $lib $fixed $LIBRARY
            install_name_tool -change $lib $fixed $LIBRARY
        fi
    done;
elif [ "x$ACTION" = "xid" ]; then
#    echo "$LIBRARY $WRONG_PREFIX to $RIGHT_PREFIX"
    lib=$(otool -D "$LIBRARY" 2>/dev/null | grep ^"$WRONG_PREFIX" | sed s,"$WRONG_PREFIX",,)
    if [ -n "$lib" ]; then
#        echo "Rewrite $lib"
        install_name_tool -id "${RIGHT_PREFIX}/${lib}" $LIBRARY;
#    else
#        path=$(otool -D "$LIBRARY" 2>/dev/null | sed -n 2p)
#        echo "Empty Result $path"
    fi
fi