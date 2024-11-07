package e2e

import "strings"

func classIssueCmd(classID, className, symbol, description, uri, data string, mintRestricted, updateRestricted bool) (command []string) {
	command = []string{
		"nft", "issue",
		classID,
	}

	if len(strings.TrimSpace(className)) > 0 {
		command = append(command, []string{"--name", className}...)
	}

	if len(strings.TrimSpace(symbol)) > 0 {
		command = append(command, []string{"--symbol", symbol}...)
	}

	if len(strings.TrimSpace(description)) > 0 {
		command = append(command, []string{"--description", description}...)
	}

	if len(strings.TrimSpace(uri)) > 0 {
		command = append(command, []string{"--uri", uri}...)
	}

	if len(strings.TrimSpace(data)) > 0 {
		command = append(command, []string{"--data", uri}...)
	}

	if mintRestricted {
		command = append(command, "--mint-restricted")
	}
	if updateRestricted {
		command = append(command, "--update-restricted")
	}
	return command
}

func tokenMintCmd(classID, tokenID, symbol, description, uri, uriHash, data string) (command []string) {
	command = []string{
		"nft", "mint",
		classID,
		tokenID,
	}

	if len(strings.TrimSpace(uri)) > 0 {
		command = append(command, []string{"--uri", uri}...)
	}

	if len(strings.TrimSpace(uriHash)) > 0 {
		command = append(command, []string{"--uri-hash", uriHash}...)
	}

	if len(strings.TrimSpace(data)) > 0 {
		command = append(command, []string{"--data", uri}...)
	}

	return command
}

func tokenInterTransferCmd(port, channel, receiver, classID string, tokenIDs ...string) (command []string) {
	command = []string{
		"nft-transfer", "transfer",
		port,
		channel,
		receiver,
		classID,
	}
	if len(tokenIDs) == 0 {
		panic("token IDs must be specified")
	}
	ids := strings.Join(tokenIDs, ",")
	command = append(command, ids)
	return command
}

func tokenQueryCmd(classID, tokenID string) (command []string) {
	command = []string{
		"nft", "token",
		classID,
		tokenID,
	}
	return command
}
