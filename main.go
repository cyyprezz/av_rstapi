package main

//C:\\privbackup.FDB"
func main() {

	a := App{}
	a.Initialize("SYSDBA", "masterkey", "C:\\privbackup.FDB")

	a.Run(":8081")
}
