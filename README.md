# owlvault


```shell
export OWLVAULT_CONFIG_PATH=/home/nitin/owlvault-go/config.yaml
```


owlvault/
├── config.yaml     # Configuration file
├── storage/        # Storage related packages
│   ├── rdbms/      # Package for RDBMS storage
│   │   ├── mysql/  # MySQL implementation
│   │   │   └── mysql_storage.go
│   │   ├── postgresql/ # PostgreSQL implementation
│   │   │   └── postgresql_storage.go
│   │   ├── mssql/  # Microsoft SQL Server implementation
│   │   │   └── mssql_storage.go
│   │   ├── oracle/ # Oracle Database implementation
│   │   │   └── oracle_storage.go
│   │   └── rdbms_storage.go # Interface for RDBMS storage
│   ├── mongodb/    # Package for MongoDB storage
│   │   └── mongodb_storage.go
│   ├── dynamodb/   # Package for DynamoDB storage
│   │   └── dynamodb_storage.go
│   └── storage.go  # Interface for storage
├── main.go         # Main entry point of the application
└── go.mod          # Go modules file