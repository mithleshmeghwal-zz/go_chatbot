
package main 


func main() {
	a := App{}
	err := a.CreateRoutes()
	if err != nil {
		panic(err)
	}
	a.Run()
}