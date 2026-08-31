CREATE USER auth_user
WITH LOGIN PASSWORD 'auth_secret';

CREATE DATABASE authdb
WITH OWNER auth_user;

CREATE USER dispatch_user
WITH LOGIN PASSWORD 'dispatch_secret';

CREATE DATABASE dispatchdb
WITH OWNER dispatch_user;

CREATE USER driver_user
WITH LOGIN PASSWORD 'driver_secret';

CREATE DATABASE driverdb
WITH OWNER driver_user;
