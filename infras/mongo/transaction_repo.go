package mongo

import (
	"context"

	"go-clean-template/entity"
	"go-clean-template/infras/mongo/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	WalletCollection        = "wallets"
	TransactionsCollection  = "transactions"
	LinkedAccountCollection = "linked_accounts"
)

type TransactionRepo struct {
	db *mongo.Database
}

func NewTransactionRepo(db *mongo.Database) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	var (
		wallet       *entity.Wallet
		walletSchema schema.WalletSchema
	)

	walletIDObj, err := primitive.ObjectIDFromHex(walletID)
	if err != nil {
		return nil, err
	}

	if err := r.db.Collection(WalletCollection).FindOne(ctx, schema.WalletSchema{ID: walletIDObj}).Decode(&walletSchema); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	wallet = walletSchema.ToWallet()
	return wallet, nil
}

func (r *TransactionRepo) SaveTransaction(ctx context.Context, trans *entity.Transaction) error {
	transSchema := schema.ToTransactionSchema(trans)

	_, err := r.db.Collection(TransactionsCollection).InsertOne(ctx, transSchema)

	return err
}

func (r *TransactionRepo) GetLinkedAccountByID(ctx context.Context, accountID string) (*entity.LinkedAccount, error) {
	var (
		account       *entity.LinkedAccount
		accountSchema schema.LinkedAccountSchema
	)
	objectId, _ := primitive.ObjectIDFromHex(accountID)
	if err := r.db.Collection(LinkedAccountCollection).FindOne(ctx, schema.LinkedAccountSchema{ID: objectId}).
		Decode(&accountSchema); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	account = accountSchema.ToLinkedAccount()
	return account, nil
}

func (r *TransactionRepo) GetBalanceByWalletID(ctx context.Context, walletID string) (float64, error) {
	var result struct {
		Balance float64 `bson:"balance"`
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"wallet_id", walletID}, {"status", entity.TransactionStatusSuccessful}}}},
		{
			{"$group", bson.D{
				{"_id", nil},
				{"balance", bson.D{
					{"$sum", bson.D{
						{"$cond", bson.D{
							{"if", bson.D{{"$eq", bson.A{"$transaction_kind", entity.TransactionIn}}}},
							{"then", "$amount"},
							{"else", bson.D{{"$multiply", bson.A{-1, "$amount"}}}},
						}},
					}},
				}},
			}},
		},
	}

	cursor, err := r.db.Collection(TransactionsCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	// Iterate over the cursor to get the result
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	// Handle errors in the aggregation cursor
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return result.Balance, nil
}

func (r *TransactionRepo) GetTransactionByID(ctx context.Context, transID string) (*entity.Transaction, error) {
	var transSchema schema.TransactionSchema

	transIdObj, err := primitive.ObjectIDFromHex(transID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Collection(TransactionsCollection).FindOne(ctx, schema.TransactionSchema{ID: transIdObj}).
		Decode(&transSchema); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return transSchema.ToTransaction(), nil
}

func (r *TransactionRepo) UpdateTransactionStatus(ctx context.Context, transID string, status entity.TransactionStatus) error {
	update := bson.D{{"$set", bson.D{{"status", string(status)}}}}
	transIDObj, err := primitive.ObjectIDFromHex(transID)
	if err != nil {
		return err
	}
	_, err = r.db.Collection(TransactionsCollection).UpdateByID(ctx, transIDObj, update)
	return err
}
