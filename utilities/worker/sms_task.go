package worker

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
//
//

// 	"github.com/hibiken/asynq"
// 	"github.com/rs/zerolog/log"

// 	"github.com/Smylet/symlet-backend/api/users"
// 	"github.com/Smylet/symlet-backend/utilities/mail"
// 	"github.com/Smylet/symlet-backend/utilities/utils"
// )

// const (

// 	TaskSendWelcomeSMS = "task:send_welcome_sms"
// )

// type PayloadSendWelcomeSMS struct {
// 	PhoneNumber string `json:"phone_number"`
// 	Message string `json:"message"`
// }

// func (distributor *RedisTaskDistributor) DistributeTaskSendWelcomeSMS(
// 	ctx context.Context,
// 	payload *PayloadSendVerifyEmail,
// 	opts ...asynq.Option,
// ) error {

// }
