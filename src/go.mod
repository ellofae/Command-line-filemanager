module main

go 1.20

replace main/wc => ./wc

replace main/cat => ./cat

require (
	main/cat v0.0.0-00010101000000-000000000000 // indirect
	main/wc v0.0.0-00010101000000-000000000000 // indirect
)
