PARAM1=$1
PARAM2=$2

# Function section
function serve_book() {
    LANG=$1

    # Rebuild english version, updating the `messages.pot` file where all chunks of texts
    # to be translated are extracted.
    MDBOOK_OUTPUT='{"xgettext": {"pot-file": "messages.pot"}}' \
    mdbook build -d po

    # Build available LANGUAGES. po file must exist for the languages listed
    # in the LANGUAGE file. If it's not, you can add a language by running (xx replaced by the language code):
    # msginit -i po/messages.pot -l xx -o po/xx.po
    for po_lang in $(cat ./LANGUAGES)
    do
        echo merging and building "$po_lang"
        msgmerge --update po/"$po_lang".po po/messages.pot
        MDBOOK_BOOK__LANGUAGE="$po_lang" mdbook build -d book/"$po_lang"
    done

    # Serving the language, if any.
    if [ -z "$LANG" ]
    then
        echo ""
        echo "No input language, stop after build."
        exit 0
    fi

    # Serve the input language, if available.
    MDBOOK_BOOK__LANGUAGE="$LANG" mdbook serve -d book/"$LANG"
}

function build_new_language() {
    LANG=$1
    FILE=po/"$LANG".po

    # Stop if no language parameter.
    if [ -z "$LANG" ]
    then
        echo ""
        echo "No input language, stop."
        exit 0
    fi

    # Build a new LANGUAGE .po file if not exist .
    if test -f "$FILE"; then
        echo ""
        echo "File $FILE already exists. Ensure you're initiating a new translation. Alternatively, use './transaction.sh $LANG' to serve an existing one."
    else
        msginit -i po/messages.pot -l "$LANG" -o po/"$LANG".po
    fi
}

# 如果参数 PARAM1 的值是 "new"，则表示要构建一个新的语言版本。此时，将 PARAM2 的值赋给变量 LANG，并调用 build_new_language 函数来构建新的语言版本。
# 如果参数 PARAM1 的值不是 "new"，则表示要构建和服务现有的语言版本。此时，将 PARAM1 的值赋给变量 LANG，并调用 serve_book 函数来构建和服务指定语言的书籍。
if [ "$PARAM1" == "new" ]; then
    # The first parameter is 'new', PARAM2 is set to $LANG.
    LANG=$PARAM2
    build_new_language "$LANG"

else
    # The first parameter is not 'new', PARAM1 is set to $LANG.
    LANG=$PARAM1
    serve_book "$LANG"

fi
