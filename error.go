package main

func CheckError(err error) {
	if err != nil {
		ErrorLogger.Println(err)
		panic(err)
	}
}
