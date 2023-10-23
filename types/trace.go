package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func ParseHexHash(hexHash string) (tmbytes.HexBytes, error) {
	if strings.TrimSpace(hexHash) == "" {
		return nil, fmt.Errorf("empty hex hash")
	}
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}

	return hash, tmtypes.ValidateHash(hash)
}

// GetClassPrefix returns the receiving class prefix
func GetClassPrefix(portID, channelID string) string {
	return fmt.Sprintf("%s/%s/", portID, channelID)
}

// RemoveClassPrefix returns the unprefixed classID.
// After the receiving chain receives the packet,if isAwayFromOrigin=false, it means that nft is moving
// in the direction of the original chain, and the portID/channelID prefix of the sending chain
// in trace.path needs to be removed
func RemoveClassPrefix(portID, channelID, classID string) (string, error) {
	classPrefix := GetClassPrefix(portID, channelID)
	if strings.HasPrefix(classID, classPrefix) {
		return strings.TrimPrefix(classID, classPrefix), nil
	}
	return "", fmt.Errorf("invalid class:%s, no class prefix: %s", classID, classPrefix)
}

// IsAwayFromOrigin determine if non-fungible token is moving away from
// the origin chain (the chain issued by the native nft).
// Note that fullClassPath refers to the full path of the unencoded classID.
// The longer the fullClassPath, the farther it is from the origin chain
func IsAwayFromOrigin(sourcePort, sourceChannel, fullClassPath string) bool {
	prefixClassID := GetClassPrefix(sourcePort, sourceChannel)
	if !strings.HasPrefix(fullClassPath, prefixClassID) {
		return true
	}
	return fullClassPath[:len(prefixClassID)] != prefixClassID
}

// ParseClassTrace parses a string with the ibc prefix (class trace) and the base classID
// into a ClassTrace type.
//
// Examples:
//
//   - "port-1/channel-1/class-1" => ClassTrace{Path: "port-1/channel-1", BaseClassId: "class-1"}
//   - "port-1/channel-1/class/1" => ClassTrace{Path: "port-1/channel-1", BaseClassId: "class/1"}
//   - "class-1" => ClassTrace{Path: "", BaseClassId: "class-1"}
func ParseClassTrace(rawClassID string) ClassTrace {
	classSplit := strings.Split(rawClassID, "/")

	if classSplit[0] == rawClassID {
		return ClassTrace{
			Path:        "",
			BaseClassId: rawClassID,
		}
	}

	path, baseClassId := extractPathAndBaseFromFullClassID(classSplit)
	return ClassTrace{
		Path:      path,
		BaseClassId: baseClassId,
	}
}

// GetFullClassPath returns the full classId according to the ICS721 specification:
// tracePath + "/" + BaseClassId
// If there exists no trace then the base BaseClassId is returned.
func (ct ClassTrace) GetFullClassPath() string {
	if ct.Path == "" {
		return ct.BaseClassId
	}
	return ct.GetPrefix() + ct.BaseClassId
}

// GetPrefix returns the receiving classId prefix composed by the trace info and a separator.
func (ct ClassTrace) GetPrefix() string {
	return ct.Path + "/"
}

// Hash returns the hex bytes of the SHA256 hash of the ClassTrace fields using the following formula:
//
// hash = sha256(tracePath + "/" + baseClassId)
func (ct ClassTrace) Hash() tmbytes.HexBytes {
	hash := sha256.Sum256([]byte(ct.GetFullClassPath()))
	return hash[:]
}

// IBCClassID a classID for an ICS721 non-fungible token in the format
// 'ibc/{hash(tracePath + BaseClassId)}'. If the trace is empty, it will return the base classID.
func (ct ClassTrace) IBCClassID() string {
	if ct.Path != "" {
		return fmt.Sprintf("%s/%s", ClassPrefix, ct.Hash())
	}
	return ct.BaseClassId
}

// Validate performs a basic validation of the ClassTrace fields.
func (ct ClassTrace) Validate() error {
	// empty trace is accepted when token lives on the original chain
	switch {
	case ct.Path == "" && ct.BaseClassId != "":
		return nil
	case strings.TrimSpace(ct.BaseClassId) == "":
		return fmt.Errorf("base class_id cannot be blank")
	}

	// NOTE: no base class validation

	identifiers := strings.Split(ct.Path, "/")
	return validateTraceIdentifiers(identifiers)
}

func validateTraceIdentifiers(identifiers []string) error {
	if len(identifiers) == 0 || len(identifiers)%2 != 0 {
		return fmt.Errorf("trace info must come in pairs of port and channel identifiers '{portID}/{channelID}', got the identifiers: %s", identifiers)
	}

	// validate correctness of port and channel identifiers
	for i := 0; i < len(identifiers); i += 2 {
		if err := host.PortIdentifierValidator(identifiers[i]); err != nil {
			return sdkerrors.Wrapf(err, "invalid port ID at position %d", i)
		}
		if err := host.ChannelIdentifierValidator(identifiers[i+1]); err != nil {
			return sdkerrors.Wrapf(err, "invalid channel ID at position %d", i)
		}
	}
	return nil
}

// Traces defines a wrapper type for a slice of DenomTrace.
type Traces []ClassTrace

// Validate performs a basic validation of each denomination trace info.
func (t Traces) Validate() error {
	seenTraces := make(map[string]bool)
	for i, trace := range t {
		hash := trace.Hash().String()
		if seenTraces[hash] {
			return fmt.Errorf("duplicated class trace with hash %s", trace.Hash())
		}

		if err := trace.Validate(); err != nil {
			return sdkerrors.Wrapf(err, "failed class trace %d validation", i)
		}
		seenTraces[hash] = true
	}
	return nil
}

var _ sort.Interface = Traces{}

// Len implements sort.Interface for Traces
func (t Traces) Len() int { return len(t) }

// Less implements sort.Interface for Traces
func (t Traces) Less(i, j int) bool { return t[i].GetFullClassPath() < t[j].GetFullClassPath() }

// Swap implements sort.Interface for Traces
func (t Traces) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Sort is a helper function to sort the set of denomination traces in-place
func (t Traces) Sort() Traces {
	sort.Sort(t)
	return t
}

// extractPathAndBaseFromFullClassID returns the trace path and the base classID from
// the elements that constitute the complete classID.
func extractPathAndBaseFromFullClassID(fullClassIdItems []string) (string, string) {
	var (
		pathSlice        []string
		baseClassIdSlice []string
	)

	length := len(fullClassIdItems)
	for i := 0; i < length; i += 2 {
		// The IBC specification does not guarantee the expected format of the
		// destination port or destination channel identifier. A short term solution
		// to determine base classID is to expect the channel identifier to be the
		// one ibc-go specifies. A longer term solution is to separate the path and base
		// denomination in the ICS721 packet.
		if i < length-1 && length > 2 && channeltypes.IsValidChannelID(fullClassIdItems[i+1]) {
			pathSlice = append(pathSlice, fullClassIdItems[i], fullClassIdItems[i+1])
		} else {
			baseClassIdSlice = fullClassIdItems[i:]
			break
		}
	}

	path := strings.Join(pathSlice, "/")
	baseClassID := strings.Join(baseClassIdSlice, "/")

	return path, baseClassID
}
