package validate

import (
	"testing"
)

type TestStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"min=18"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			s: TestStruct{
				Name: "John Doe",
				Age:  25,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			s: TestStruct{
				Age: 25,
			},
			wantErr: true,
		},
		{
			name: "age too young",
			s: TestStruct{
				Name: "John Doe",
				Age:  17,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
