package user

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoginCheckerService_Check(t *testing.T) {
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "success",
			args: args{
				login: "test",
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "short",
			args: args{
				login: "t",
			},
			wantErr: true,
			err:     ErrLoginNotValid,
		},
		{
			name: "long",
			args: args{
				login: strings.Repeat("a", 129),
			},
			wantErr: true,
			err:     ErrLoginNotValid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewLoginCheckerService()
			err := checker.Check(tt.args.login)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
