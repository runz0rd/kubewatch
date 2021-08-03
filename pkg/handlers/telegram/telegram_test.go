/*
Copyright 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package telegram

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestTelegramInit(t *testing.T) {
	s := &Telegram{}
	expectedErrorBotToken := fmt.Errorf(telegramErrMsg, "Missing Telegram bot token")
	expectedErrorChatId := fmt.Errorf(telegramErrMsg, "Missing Telegram chat id")

	var Tests = []struct {
		telegram config.Telegram
		err      error
	}{
		{config.Telegram{}, expectedErrorBotToken},
		{config.Telegram{BotToken: "test"}, expectedErrorChatId},
		{config.Telegram{BotToken: "", ChatId: 0}, nil},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Telegram = tt.telegram
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
