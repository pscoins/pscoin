cockroach user set pscoinroach --insecure
cockroach sql --insecure -e 'CREATE DATABASE pscoin'
cockroach sql --insecure -e 'GRANT ALL ON DATABASE pscoin TO pscoinroach'
