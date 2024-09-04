package cmd

import (
	"fmt"
	"os"
	"taskmanager/internal/tasks"
	UIForm "taskmanager/ui/form"
	UITable "taskmanager/ui/table"
	"taskmanager/utils"
	"time"

	// "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

const VERSION = "1.0.5"

func Execute(db *gorm.DB) {

	var cmdAdd = &cobra.Command{
		Use:   "add",
		Short: "Command to add a task",
		Long:  `This command add a task to the task list.`,
		Run: func(cmd *cobra.Command, args []string) {
			task := UIForm.Create()
			newTask := tasks.Add(db, task)

			fmt.Printf("Tarea creada con éxito: %s(%d)\n", task.Name, newTask.ID)
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Command to list all tasks",
		Long:  "This command list all tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			tasksList := tasks.GetAll(db)
			if len(tasksList) == 0 {
				fmt.Println("There are no tasks to do")
				return
			}

			/* columns := []table.Column{
				{Title: "ID", Width: 4},
				{Title: "Nombre", Width: 10},
				{Title: "Descripción", Width: 20},
				{Title: "Estado", Width: 10},
			} */

			values := []table.Row{}
			for _, task := range tasksList {
				done := "❌" // "✅"
				incompleteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#af00d7"))
				completeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00d700")).Bold(true)
				var row table.Row

				if task.Completed {
					done = "✅"
					row = table.NewRow(table.RowData{
						UITable.ColumnKeyID:          task.ID,
						UITable.ColumnKeyName:        task.Name,
						UITable.ColumnKeyDescription: task.Description,
						UITable.ColumnKeyStatus:      done,
					}).WithStyle(completeStyle)
				} else {
					row = table.NewRow(table.RowData{
						UITable.ColumnKeyID:          task.ID,
						UITable.ColumnKeyName:        task.Name,
						UITable.ColumnKeyDescription: task.Description,
						UITable.ColumnKeyStatus:      done,
					}).WithStyle(incompleteStyle)
				}

				values = append(values, row)
			}

			// for i := 0; i < len(tasksList); i++ {
			// 	values = append(values, table.Row{strconv.Itoa(tasksList[i].ID), tasksList[i].Name, tasksList[i].Description, strconv.FormatBool(tasksList[i].Completed)})
			// }

			UITable.NewModel(values)
		},
	}

	var cmdDetail = &cobra.Command{
		Use:   "detail [id]",
		Short: "Command to show a task detail",
		Long:  `This command show a task detail.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := utils.ParseInt(args[0])
			task := tasks.GetByID(db, id)
			fmt.Printf("ID: %d \nName: %s \nDescription: %s \nStatus: %t \nCreated At: %s \n", task.ID, task.Name, task.Description, task.Completed, task.CreatedAt.Format(time.RFC822Z))

		},
	}

	var cmdUpdate = &cobra.Command{
		Use:   "update [id] [string] [string]",
		Short: "Command to update a task",
		Long:  `This command update a task.`,
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := utils.ParseInt(args[0])

			var tasksEdited = tasks.GetByID(db, id)

			if args[1] == "name" {
				tasksEdited.Name = args[2]
			} else if args[1] == "description" {
				tasksEdited.Description = args[2]
			} else {
				panic("invalid argument: some field is not valid")
			}

			task := tasks.UpdateByID(db, id, *tasksEdited)

			fmt.Printf("Actualización realizada con éxito \n\nID: %d \nNombre: %s \nDescripción: %s \nEstado: %t \n", task.ID, task.Name, task.Description, task.Completed)

		},
	}

	var cmdToggled = &cobra.Command{
		Use:   "toggled [id]",
		Short: "Command to mark a task as completed/incomplete",
		Long:  `This command mark a task as completed/incomplete.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := utils.ParseInt(args[0])

			var tasksEdited = tasks.GetByID(db, id)

			tasksEdited.Completed = !tasksEdited.Completed

			task := tasks.UpdateByID(db, id, *tasksEdited)

			fmt.Printf("La tarea con ID: %d ha sido cambiada de status\n", task.ID)
		},
	}

	var cmdDeleted = &cobra.Command{
		Use:   "delete [id]",
		Short: "Command to delete a task",
		Long:  `This command delete a task.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := utils.ParseInt(args[0])

			tasks.DeleteByID(db, id)

			fmt.Printf("La tarea con ID: %d se ha eliminado\n", id)
		},
	}

	var rootCmd = &cobra.Command{
		Use:     "task",
		Version: VERSION,
		Short:   "task is a command line application for managing to-do tasks.",
	}
	rootCmd.AddCommand(cmdAdd)
	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdDetail)
	rootCmd.AddCommand(cmdUpdate)
	rootCmd.AddCommand(cmdToggled)
	rootCmd.AddCommand(cmdDeleted)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/* REFERENCES:
https://sweworld.net/cheatsheets/terminal_escape_code/
https://github.com/Evertras/bubble-table/issues/179
https://github.com/Evertras/bubble-table/tree/main/examples/features
https://www.freedium.cfd/https://dlcoder.medium.com/lipgloss-in-golang-beautifully-style-your-command-line-applications-d14fda470906

*/
