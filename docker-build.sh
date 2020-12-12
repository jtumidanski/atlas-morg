if [ $1 = NO-CACHE ]
then
   docker build --no-cache --tag atlas-morg:latest .
else
   docker build --tag atlas-morg:latest .
fi
