@ECHO OFF
IF "%1"=="NO-CACHE" docker build -f Dockerfile.dev --no-cache --tag atlas-morg:latest .
IF NOT "%1"=="NO-CACHE" docker -f Dockerfile.dev build --tag atlas-morg:latest .
