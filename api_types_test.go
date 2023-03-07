// Файл api_types.go содержит описания структур и их валидацию на уровне API сервера

package main

import (
	"testing"
	"time"
)

func Test_checkAgeOver18(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "born after today",
			wantErr: true,
			args:    args{value: time.Date(2038, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:    "born yesterday",
			wantErr: true,
			args:    args{value: time.Now().Add(-time.Hour * 24)},
		},
		{
			name:    "born in 90s",
			wantErr: false,
			args:    args{value: time.Date(1995, 2, 2, 0, 0, 0, 0, time.UTC)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAgeOver18(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("checkAgeOver18() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
