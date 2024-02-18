
## services PgDatabase


#### Config services

## SCHEMENAME: USERS, ...
# PORT: 3022
# PATH_NAME: C:\Users\PC\Documents\PgDatabase


+ Description:

### servicio db el cual permite la gestion y comunicacion de bases de datos con los microservicios



+ Services:

* PATH: /:schema/:table     METHOD: GET  
* PATH: /:schema/:table     METHOD: POST  
* PATH: /:schema/:table     METHOD: DELETE 
* PATH: /:schema/:table     METHOD: PUT
* PATH: /:schema/:table     METHOD: MERGE   


+ ERROR STATUS CODE 

# ERR: 500
