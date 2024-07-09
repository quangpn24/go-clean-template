package usecase

import (
	"go-clean-template/usecase/interfaces"
	"go-clean-template/usecase/mocks"
	"reflect"
	"testing"
)

func TestNewTransactionUseCase(t *testing.T) {
	type args struct {
		repo          interfaces.ITransactionRepository
		bankSvc       interfaces.IBankService
		dbTransaction interfaces.IDBTransaction
	}
	tests := []struct {
		name string
		args args
		want *TransactionUseCase
	}{
		{
			name: "create new transaction use case",
			args: args{
				repo:          mocks.NewITransactionRepository(t),
				bankSvc:       mocks.NewIBankService(t),
				dbTransaction: mocks.NewIDBTransaction(t),
			},
			want: &TransactionUseCase{
				repo:          mocks.NewITransactionRepository(t),
				bankSvc:       mocks.NewIBankService(t),
				notifiers:     []interfaces.INotifier{},
				dbTransaction: mocks.NewIDBTransaction(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionUseCase(tt.args.repo, tt.args.bankSvc, tt.args.dbTransaction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
