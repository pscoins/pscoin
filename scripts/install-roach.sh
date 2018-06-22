wget -qO- https://binaries.cockroachdb.com/cockroach-v2.0.3.linux-amd64.tgz | tar  xvz
cp -i cockroach-v2.0.3.linux-amd64/cockroach /usr/local/bin

cockroach user set pscoinroach --insecure
cockroach sql --insecure -e 'CREATE DATABASE pscoin'
cockroach sql --insecure -e 'GRANT ALL ON DATABASE pscoin TO pscoinroach'

