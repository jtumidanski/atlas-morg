@ECHO OFF
IF "%1"=="NO-CACHE" docker build --no-cache --tag atlas-morg:latest .
IF NOT "%1"=="NO-CACHE" docker build --tag atlas-morg:latest .