package matrix

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/commitments.json
var refTests []byte

func TestStateCommitment(t *testing.T) {
	tests := []struct {
		expectedPacked string
		matrix         []uint64 // Automatically padded with 0s to the required length
	}{
		{
			expectedPacked: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			expectedPacked: "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000700000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000d000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000f0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000013000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000150000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001700000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000000000000000019",
			matrix:         []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		},
		{
			expectedPacked: "000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			matrix:         []uint64{18446744073709551615},
		},
	}
	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			state := NewStateMatrix()
			copy(state.s.a[:], test.matrix)
			expected := crypto.Keccak256Hash(common.Hex2Bytes(test.expectedPacked))
			actual := state.StateCommitment()
			require.Equal(t, test.expectedPacked, common.Bytes2Hex(state.PackState()))
			require.Equal(t, expected, actual)
		})
	}
}

type testData struct {
	Input       []byte        `json:"input"`
	Commitments []common.Hash `json:"commitments"`
}

func TestReferenceCommitments(t *testing.T) {
	var tests []testData
	require.NoError(t, json.Unmarshal(refTests, &tests))

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Ref-%v", i), func(t *testing.T) {
			s := NewStateMatrix()
			commitments := []common.Hash{s.StateCommitment()}
			for i := 0; i < len(test.Input); i += LeafSize {
				end := min(i+LeafSize, len(test.Input))
				s.AbsorbLeaf(test.Input[i:end], end == len(test.Input))
				commitments = append(commitments, s.StateCommitment())
			}
			if len(test.Input) == 0 {
				s.AbsorbLeaf(nil, true)
				commitments = append(commitments, s.StateCommitment())
			}
			actual := s.Hash()
			expected := crypto.Keccak256Hash(test.Input)
			require.Equal(t, expected, actual)
			require.Equal(t, test.Commitments, commitments)
		})
	}
}

func TestReferenceCommitmentsFromReader(t *testing.T) {
	var tests []testData
	require.NoError(t, json.Unmarshal(refTests, &tests))

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Ref-%v", i), func(t *testing.T) {
			s := NewStateMatrix()
			commitments := []common.Hash{s.StateCommitment()}
			in := bytes.NewReader(test.Input)
			for {
				err := s.AbsorbNextLeaf(in)
				if errors.Is(err, io.EOF) {
					commitments = append(commitments, s.StateCommitment())
					break
				}
				// Shouldn't get any error except EOF
				require.NoError(t, err)
				commitments = append(commitments, s.StateCommitment())
			}
			actual := s.Hash()
			expected := crypto.Keccak256Hash(test.Input)
			require.Equal(t, expected, actual)
			require.Equal(t, test.Commitments, commitments)
		})
	}
}

func FuzzKeccak(f *testing.F) {
	f.Fuzz(func(t *testing.T, number, time uint64, data []byte) {
		s := NewStateMatrix()
		for i := 0; i < len(data); i += LeafSize {
			end := min(i+LeafSize, len(data))
			s.AbsorbLeaf(data[i:end], end == len(data))
		}
		actual := s.Hash()
		expected := crypto.Keccak256Hash(data)
		require.Equal(t, expected, actual)
	})
}
