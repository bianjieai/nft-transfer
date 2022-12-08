package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonFungibleTokenPacketData_ValidateBasic(t *testing.T) {
	tokenData := []byte{}
	tests := []struct {
		name    string
		packet  NonFungibleTokenPacketData
		wantErr bool
	}{
		{
			name:    "valid packet",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{"kitty_uri"}, [][]byte{tokenData}, sender, receiver},
			wantErr: false,
		},
		{
			name:    "invalid packet with empty classID",
			packet:  NonFungibleTokenPacketData{"", "uri", []string{"kitty"}, []string{"kitty_uri"}, [][]byte{tokenData}, sender, receiver},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty tokenIds",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{}, []string{"kitty_uri"}, [][]byte{tokenData}, sender, receiver},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty tokenUris",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, [][]byte{tokenData}, sender, receiver},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty sender",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, [][]byte{tokenData}, "", receiver},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty receiver",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, [][]byte{tokenData}, sender, receiver},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.packet.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("NonFungibleTokenPacketData.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNonFungibleTokenPacketData_Optimize(t *testing.T) {
	tokenData := []byte{}
	tests := []struct {
		name   string
		packet NonFungibleTokenPacketData
		len    int
	}{
		{
			name:   "empty tokenData",
			packet: NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{"kitty_uri"}, [][]byte{tokenData}, sender, receiver},
			len:    0,
		},
		{
			name:   "one tokenData",
			packet: NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{"kitty_uri"}, [][]byte{[]byte("tokenData")}, sender, receiver},
			len:    1,
		},
		{
			name:   "two tokenData",
			packet: NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty", "kitty1"}, []string{"kitty_uri", "kitty_uri1"}, [][]byte{[]byte("tokenData"), []byte("tokenData1")}, sender, receiver},
			len:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := tt.packet.Optimize()
			require.Equal(t, len(t1.TokenData), tt.len)
		})
	}
}
