package simapp

// import (
// 	"github.com/bianjieai/nft-transfer/testing/mock"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// )

// // Setup initializes a new SimApp. A Nop logger is set in SimApp.
// func Setup(isCheckTx bool) *SimApp {
// 	privVal := mock.NewPV()
// 	pubKey, _ := privVal.GetPubKey()

// 	// create validator set with single validator
// 	validator := tmtypes.NewValidator(pubKey, 1)
// 	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

// 	// generate genesis account
// 	senderPrivKey := secp256k1.GenPrivKey()
// 	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
// 	balance := banktypes.Balance{
// 		Address: acc.GetAddress().String(),
// 		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
// 	}

// 	app := SetupWithGenesisValSet(valSet, []authtypes.GenesisAccount{acc}, balance)

// 	return app
// }
