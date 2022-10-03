
go build -o booking cmd/web/*.go
./booking -dbname=bookings -dbuser=postgres -dbpass=1234 -cache=false -production=false