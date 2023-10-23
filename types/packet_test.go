package types

import (
	"bytes"
	"testing"
)

func TestNonFungibleTokenPacketData_ValidateBasic(t *testing.T) {
	tokenData := []string{}
	tests := []struct {
		name    string
		packet  NonFungibleTokenPacketData
		wantErr bool
	}{
		{
			name:    "valid packet",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{"kitty_uri"}, tokenData, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "invalid packet with empty classID",
			packet:  NonFungibleTokenPacketData{"", "uri", "", []string{"kitty"}, []string{"kitty_uri"}, tokenData, sender, receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty tokenIds",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{}, []string{"kitty_uri"}, tokenData, sender, receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "invalid packet with repeated tokenIds",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty","kitty"}, []string{"kitty_uri","kitty_uri"}, tokenData, sender, receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "valid packet with empty tokenUris",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, tokenData, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with nil tokenUris",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, nil, tokenData, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with tokenUris",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{"1"}, tokenData, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with tokenUris of empty string entry",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty", "mary"}, []string{"1", ""}, tokenData, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "invalid packet with unmatched tokenUris number",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{"1", "2"}, tokenData, sender, receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "valid packet with empty tokenData",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, []string{}, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with nil tokenData",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, nil, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with tokenData",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, []string{"1"}, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "valid packet with tokenData of empty string entry",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty", "mary"}, []string{}, []string{"1", ""}, sender, receiver, "memo"},
			wantErr: false,
		},
		{
			name:    "invalid packet with unmatched tokenData number",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, []string{"1", "2"}, sender, receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty sender",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, tokenData, "", receiver, "memo"},
			wantErr: true,
		},
		{
			name:    "invalid packet with empty receiver",
			packet:  NonFungibleTokenPacketData{"cryptoCat", "uri", "", []string{"kitty"}, []string{}, tokenData, sender, "", "memo"},
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

func TestNonFungibleTokenPacketData_GetBytes(t *testing.T) {
	type fields struct {
		ClassId   string
		ClassUri  string
		ClassData string
		TokenIds  []string
		TokenUris []string
		TokenData []string
		Sender    string
		Receiver  string
		Memo      string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"success",
			fields{"classId", "classUri", "classData", []string{"id1", "id2"}, []string{"uri1", "uri2"}, []string{"data1"}, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"],"tokenUris":["uri1","uri2"]}`),
		},
		{
			"success with missing classUri",
			fields{"classId", "", "classData", []string{"id1", "id2"}, []string{"uri1", "uri2"}, []string{"data1"}, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"],"tokenUris":["uri1","uri2"]}`),
		},
		{
			"success with missing classData",
			fields{"classId", "classUri", "", []string{"id1", "id2"}, []string{"uri1", "uri2"}, []string{"data1"}, sender, receiver, "memo"},
			[]byte(`{"classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"],"tokenUris":["uri1","uri2"]}`),
		},
		{
			"success with empty tokenUris",
			fields{"classId", "classUri", "classData", []string{"id1", "id2"}, []string{"", ""}, []string{"data1"}, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"`),
		},
		{
			"success with nil tokenUris",
			fields{"classId", "classUri", "classData", []string{"id1", "id2"}, nil, []string{"data1"}, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"`),
		},
		{
			"success with empty tokenData",
			fields{"classId", "classUri", "classData", []string{"id1", "id2"}, []string{"uri1", "uri2"}, []string{"", ""}, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"`),
		},
		{
			"success with nil tokenData",
			fields{"classId", "classUri", "classData", []string{"id1", "id2"}, []string{"uri1", "uri2"}, nil, sender, receiver, "memo"},
			[]byte(`{"classData":"classData","classId":"classId","classUri":"classUri","memo":"memo","receiver":"cosmos15mn87gny58ptfpzq0du6t398gle50xphkw2pkt","sender":"cosmos1eshqg3adwvnuqng0eqfr6ppj35j9hh6zyd9qss","tokenData":["data1"],"tokenIds":["id1","id2"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nftpd := NonFungibleTokenPacketData{
				ClassId:   tt.fields.ClassId,
				ClassUri:  tt.fields.ClassUri,
				ClassData: tt.fields.ClassData,
				TokenIds:  tt.fields.TokenIds,
				TokenUris: tt.fields.TokenUris,
				TokenData: tt.fields.TokenData,
				Sender:    tt.fields.Sender,
				Receiver:  tt.fields.Receiver,
				Memo:      tt.fields.Memo,
			}
			t.Logf("%s\n", nftpd.GetBytes())
			if got := nftpd.GetBytes(); bytes.Equal(got, tt.want) {
				t.Errorf("NonFungibleTokenPacketData.GetBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
