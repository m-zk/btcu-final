package main

import (
	"btcu-final/client"
	"log"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	println("GenerateKeys")
	pri, pub, err := client.GenerateKeys()
	check(err)
	println(*pri,*pub)

	println("Register")
	const TEST_HASH = "12345678912345678912345679"
	data,err := client.Register(&client.Copyright{
		Name:      "TEST",
		Author:    "TEST_USER",
		Press:     "TEST_PRESS",
		Hash:      TEST_HASH,
		PublicKey: *pub,
	}, *pri)
	check(err)
	println(data)

	println("GetRightByHash")
	data, err = client.GetRightByHash(TEST_HASH)
	check(err)
	println(data)

	println("GetRightByInfo")
	data, err = client.GetRightByInfo("TEST", "TEST_USER", "TEST_PRESS")
	check(err)
	println(data)

	println("===TEST PASS===")
}
