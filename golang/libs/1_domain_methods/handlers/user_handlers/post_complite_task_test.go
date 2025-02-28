package user_handlers

import (
	"de-net/libs/1_domain_methods/helpers"
	"de-net/libs/2_generated_models/model"
	"de-net/libs/3_infrastructure/test_sqlite"
	"de-net/libs/4_common/smart_context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CompleteTask(t *testing.T) {
	db := test_sqlite.CreateUniqueSqliteInMemoryForUnitTests()
	logger := smart_context.NewSmartContext().WithDB(db)

	// Создаём пользователей и задания
	userID := helpers.GenerateUUID()
	taskID := helpers.GenerateUUID()

	db.AutoMigrate(&model.User{}, &model.Task{}, &model.Point{})
	db.Create(&model.User{ID: userID, Username: "test_user"})
	db.Create(&model.Task{ID: taskID, Name: "Test Task"})

	tests := []struct {
		name    string
		userID  string
		taskID  string
		wantErr bool
	}{
		{
			name:    "Task completion success",
			userID:  userID,
			taskID:  taskID,
			wantErr: false,
		},
		{
			name:    "Task not found",
			userID:  userID,
			taskID:  helpers.GenerateUUID(), // Несуществующий ID
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := completeTask(logger, tt.userID, tt.taskID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
