if [[ "$1" = "NO-CACHE" ]]
then
   docker build -f Dockerfile.dev --no-cache --tag atlas-morg:latest .
else
   docker build -f Dockerfile.dev --tag atlas-morg:latest .
fi
