#!/bin/sh

git filter-branch -f --env-filter '
an="$GIT_AUTHOR_NAME"
am="$GIT_AUTHOR_EMAIL"
cn="$GIT_COMMITTER_NAME"
cm="$GIT_COMMITTER_EMAIL"
if [ "$GIT_COMMITTER_EMAIL" = "you@example.com" ]
then
    cn="hunterhug"
    cm="gdccmcm14@live.com"
fi
if [ "$GIT_AUTHOR_EMAIL" = "you@example.com" ]
then
    an="hunterhug"
    am="gdccmcm14@live.com"
fi
    export GIT_AUTHOR_NAME="$an"
    export GIT_AUTHOR_EMAIL="$am"
    export GIT_COMMITTER_NAME="$cn"
    export GIT_COMMITTER_EMAIL="$cm"
if [ "$GIT_COMMITTER_EMAIL" = "gdccmcm@live.com" ]
then
    cn="hunterhug"
    cm="gdccmcm14@live.com"
fi
if [ "$GIT_AUTHOR_EMAIL" = "gdccmcm@live.com" ]
then
    an="hunterhug"
    am="gdccmcm14@live.com"
fi
    export GIT_AUTHOR_NAME="$an"
    export GIT_AUTHOR_EMAIL="$am"
    export GIT_COMMITTER_NAME="$cn"
    export GIT_COMMITTER_EMAIL="$cm"
'