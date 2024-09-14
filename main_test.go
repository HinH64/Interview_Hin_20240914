package main

import (
	"os"
	"testing"

	"Interview_Hin_20240914/services"
)

func TestMain(m *testing.M) {
    // Setup code
    services.LoadConfig()
    services.InitMongoDB()

    // Run tests
    code := m.Run()

    // Teardown code (if needed)

    // Exit with the test status code
    os.Exit(code)
}