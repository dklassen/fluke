#! /usr/bin/env bash

cd ./grammar
antlr -Dlanguage=Go -package parser -o ../parser ClickHouseLexer.g4
antlr -Dlanguage=Go -package parser -o ../parser ClickHouseParser.g4
