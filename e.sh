function e() {
    res=$(./e "$@")
    if [[ $res == \#env* ]];
    then
        source <(echo "$res")
    else
        echo "$res"
    fi
}
