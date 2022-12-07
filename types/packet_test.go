package types

import "testing"

func TestNonFungibleTokenPacketData_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		packet  NonFungibleTokenPacketData
		wantErr bool
	}{
		{
			name:    "valid packet",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{"kitty_uri"}, sender, receiver, [][]byte{nil}},
			wantErr: false,
		},
		{
			name:    "invalid packet with empty classID",
			packet:  NonFungibleTokenPacketData{"", "uri", []string{"kitty"}, []string{"kitty_uri"}, sender, receiver, [][]byte{nil}},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty tokenIds",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{}, []string{"kitty_uri"}, sender, receiver, [][]byte{nil}},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty tokenUris",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, sender, receiver, [][]byte{nil}},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty sender",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, "", receiver, [][]byte{nil}},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty receiver",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", []string{"kitty"}, []string{}, sender, receiver, [][]byte{nil}},
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
