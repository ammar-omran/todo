package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	pbmigrations "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	pbmigrations.Register(func(app core.App) error {
		collection := core.NewBaseCollection("todos")
		collection.Fields.Add(
			&core.TextField{
				Name:     "description",
				Required: true,
				Max:      5000,
			},
			&core.BoolField{
				Name: "done",
			},
			&core.AutodateField{
				Name:     "created",
				System:   true,
				OnCreate: true,
			},
			&core.AutodateField{
				Name:     "updated",
				System:   true,
				OnCreate: true,
				OnUpdate: true,
			},
		)

		collection.ListRule = types.Pointer("")
		collection.ViewRule = types.Pointer("")
		collection.CreateRule = types.Pointer("")
		collection.UpdateRule = types.Pointer("")
		collection.DeleteRule = types.Pointer("")

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("todos")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
