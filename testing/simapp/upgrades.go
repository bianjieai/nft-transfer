package simapp

import "github.com/bianjieai/nft-transfer/testing/simapp/upgrades"

// registerUpgradeHandlers registers the migration handler for the test application.
// A deployed chain must additionally supply its own store upgrades and any
// application-specific migrations, including settlement of legacy IBC fees.
func (app *SimApp) registerUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgrades.V10,
		upgrades.CreateDefaultUpgradeHandler(app.ModuleManager, app.configurator),
	)
}
