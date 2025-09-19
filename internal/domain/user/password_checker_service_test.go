package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordChecker_CheckPassword(t *testing.T) {
	type args struct {
		password string
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "strong",
			args: args{
				password: "UgL.k.0`4c1y",
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "short",
			args: args{
				password: "Ug",
			},
			wantErr: true,
			err:     ErrPasswordTooShort,
		},
		{
			name: "no-upper",
			args: args{
				password: "password1",
			},
			wantErr: true,
			err:     ErrPasswordHasNoUpper,
		},
		{
			name: "no-lower",
			args: args{
				password: "PASSWORD1",
			},
			wantErr: true,
			err:     ErrPasswordHasNoLower,
		},
		{
			name: "no-digits",
			args: args{
				password: "Password",
			},
			wantErr: true,
			err:     ErrPasswordHasNoDigit,
		},
		{
			name: "no-special",
			args: args{
				password: "PaSSsword123",
			},
			wantErr: true,
			err:     ErrPasswordHasNoSpecial,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewPasswordCheckerService()

			err := checker.Check(tt.args.password)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
