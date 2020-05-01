module eh-sender-receiver/cli/eht

go 1.14

replace eh-sender-receiver/sender => ./../../sender
replace eh-sender-receiver/receiver => ./../../receiver

replace eh-sender-receiver/cli/sender => ./../sender
replace eh-sender-receiver/cli/receiver => ./../receiver

require (
	eh-sender-receiver/receiver v0.0.0-00010101000000-000000000000 // indirect
	eh-sender-receiver/sender v0.0.0-00010101000000-000000000000 // indirect
    eh-sender-receiver/cli/receiver v0.0.0-00010101000000-000000000000 // indirect
    eh-sender-receiver/cli/sender v0.0.0-00010101000000-000000000000 // indirect
)
