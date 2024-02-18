{

     "ConfigDb":{
        "Port": ":3022",
        "Protocolo":"http",
        "Host":"localhost",
        "name":"db"
    },

    "Log":{
        "FileName":"log/msg.log",
        "console_level":"debug"
    },
    "requestBodyLimit": "50mb",
    "requestParamLimit": "50mb",
    "requestTimeout": 30000,
    "responseTimeout": 30000,
    "requestSizeLimit": "50mb",
    "responseSizeLimit": "50mb",
    "requestHeaderLimit": "50mb",
    "db_config":{
        "TimeZone": "America/UTC",
        "Lang":"LATIN AMERICA SPANISH_COLOMBIA.1252"
    },

    "Schemas":{
        "audits":"audits"
    },

    "Db_list_conn":[
        {
            "host":"monorail.proxy.rlwy.net",
            "user":"postgres",
            "password":"CbF6AFD33G*dAG3ff1c4*Bde4dFcFB53",
            "databases":"railway",
            "schema":"users",
            "port":"44896",
            "sslmode":"disable",
            "TimeZone": "America/UTC",
            "Lang":"LATIN AMERICA SPANISH_COLOMBIA.1252",
            "requestBodyLimit": "50mb",
            "requestParamLimit": "50mb",
            "requestTimeout": 30000,
            "responseTimeout": 30000

        },
         {
            "host":"localhost",
            "user":"postgres",
            "password":"1008",
            "databases":"postgres",
            "schema":"users",
            "port":"5432",
            "sslmode":"disable",
            "TimeZone": "America/UTC",
            "Lang":"LATIN AMERICA SPANISH_COLOMBIA.1252",
            "requestBodyLimit": "50mb",
            "requestParamLimit": "50mb",
            "requestTimeout": 30000,
            "responseTimeout": 30000

        }
        
        
        
    ]
}