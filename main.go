package main

// Hier den Pfad zur Datenbank hinterlegen
// Sobald der Pfad hinterlegt ist
// kann man die Applikation via go build kompilieren
// anschließend kann man die Anwendung starten z.B ./rest-api
//Author Nicklas Desens
func main() {

	a := App{}
	a.Initialize("SYSDBA", "masterkey", "C:\\privbackup.FDB")

	a.Run(":8081")
}
