package main

import "net/http"

// healthcheckHandler sends a json response representing the status
// of the application.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{}
	response["environment"] = "development"
	response["version"] = version
	err := app.writeJSON(w, r, http.StatusOK, response, nil)
	if err != nil {
		app.logger.Println(err)
	}
}
