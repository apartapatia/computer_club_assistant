package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apartapatia/computer_club_assistant/internal/myparser"
	"github.com/apartapatia/computer_club_assistant/pkg/client"
	"github.com/apartapatia/computer_club_assistant/pkg/handlers"
	"github.com/apartapatia/computer_club_assistant/pkg/table"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: computer_club_assistant <file_name>")
		fmt.Println("🪟 For Windows: ./computer_club_assistant.exe <file_name>")
		fmt.Println("🐧 For Linux: ./computer_club_assistant <file_name>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	filePath := filepath.Join("configs/" + fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File %s not found. Please check the file path and try again.\n", filePath)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	pars := myparser.NewFileParser(file)
	clubInfo, err := pars.ReadClubInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	managerInfo, err := pars.ReadManagerEvents(clubInfo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clients := client.NewMemoryRepo()
	tables := table.NewMemoryRepo(clubInfo.MaxTables)
	handler := handlers.NewCommandHandler(clubInfo, managerInfo, clients, tables)

	res := handler.HandleCommands()
	fmt.Println(res)
}
