package main

func main() {
	app, err := newApp()
	if err != nil {
		panic(err)
	}
	if err := app.run(); err != nil {
		panic(err)
	}
}
