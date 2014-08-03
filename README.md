[![Build Status](https://travis-ci.org/andrascsibi/how-many-pullups.svg?branch=master)](https://travis-ci.org/andrascsibi/how-many-pullups)

# PullApp

PullApp is for counting sets and reps. It runs on AppEngine. Here:

[http://pull-app.appspot.com](http://pull-app.appspot.com)

## Project setup

check out this project in

    $GOPATH/src/github.com/andrascsibi

make a symlink to the pre-commit hook

    ln -s ../../pre-commit.sh .git/hooks/pre-commit

get the [AppEngine SDK for Go from here](https://developers.google.com/appengine/downloads)

run it locally

    cd $GOPATH/src/github.com/andrascsibi/how-many-pullups
    goapp serve