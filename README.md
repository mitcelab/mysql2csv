# mysql2csv
It's surprisingly hard to export a MySQL query to a CSV file. You can use `SELECT ... INTO OUTFILE ...` but 
that requires special file permissions; on the other hand, third-party clients such as MySQL Workbench and 
SequelPro are slow and often run into encoding problems.

#### install
> go install github.com/mitcelab/mysql2csv

#### usage
> mysql2csv --user X --pass Y --host Z --dbname btctalk --query "SELECT * FROM User"

By default, it prints to STDOUT where you can easily pipe it through gzip/bzip2, but you can also specify 
a output file with the `--output` option.
