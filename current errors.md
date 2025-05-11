root@RWTChoir:~/RWTProj# ./force-cleanup.sh
===== PERFORMING AGGRESSIVE CONTAINER CLEANUP =====
Stopping all containers...
Stopping rwtproj-db ... done
Removing 13cd1d09efca_rwtproj-api ... done
Removing rwtproj-db               ... done
Removing network rwtproj_default
Removing volume rwtproj_db-data
Force removing all project containers...
Pruning volumes...
Total reclaimed space: 0B
Setting up SQL directory...
Copied rwtchoir.sql to sql directory
Creating necessary directories...
===== CLEANUP COMPLETE =====
You can now try restarting with just the database:
docker-compose up -d db
Then check it with: docker-compose ps
Once the database is healthy, start the API with: docker-compose up -d api
root@RWTChoir:~/RWTProj# docker-compose pull
Pulling db       ... done
Pulling api      ... done
Pulling frontend ... done
root@RWTChoir:~/RWTProj# docker-compose up -d
Creating network "rwtproj_default" with the default driver
Creating volume "rwtproj_db-data" with default driver
Creating rwtproj-db ... done
Creating rwtproj-api ... done
Creating rwtproj-frontend ... done